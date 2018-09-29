# JSL

JSL (short for the **J**ansky **S**tack **L**angauge) is a dynamically-typed and garbage collected stack-based language. It is currently a work in progress.

## Starting the JSL Interpreter

To use JSL, you can start the interpreter by running `go run *.go` in the directory containing the JSL code.

## JSL By Example

The following short examples illustrate the design and features of the JSL language.

### Arithmetic

JSL can function as a glorified interpreter

    > 3 4 +
    7.000000
    > 2 *
    14.000000
    > 3 /
    4.666667

JSL only has one numeric type, representing both integers and floating point numbers.

### Variables

Variables can be assigned to using `asn`.

    > 3.14 'pi asn
    > pi
    3.140000

The single quote that prefixes the variable name `pi` indicates that it should be treated as a reference. If a variable with the name `pi` does not exist within the current scope, then the variable name is pushed to the stack as a *name reference*.

    > 'pi
    'pi

However, if a variable `pi` does exist, then a *variable reference* to that variable will be placed onto the stack.

    > 3.14 'pi asn
    > 'pi
    <Reference: 4D65822107FCFD52>

Here, `4D65822107FCFD52` is the unique ID that the variable `pi` has on the global symbol table. References function somewhat like pointers in C. We can use the `@` operator to dereference a reference.

    > 'pi@
    3.140000

Assigning to a variable requires either a name reference or a variable reference.

    > 3.14 'pi asn
    > pi     
    3.140000
    > 3.141 'pi asn
    3.140000
    > pi
    3.141000
    3.140000
    > 3.1415 pi asn
    Error: Expected, but did not receive an identifier reference or reference.

Notice that after each command is executed, the entire contents of the stack are printed. You can run the command `clear` to clear the stack at any time.

You can assign any language object to a variable, including a reference.

    > 3.14 'pi asn
    > 'pi 'pi_ref asn
    > pi_ref
    <Reference: 4D65822107FCFD52>
    > pi_ref@
    3.140000
    <Reference: 4D65822107FCFD52>
    > 3.141 pi_ref asn
    3.140000
    <Reference: 4D65822107FCFD52>
    > pi
    3.141000
    3.140000
    <Reference: 4D65822107FCFD52>
    > 'pi_ref 'pi_ref_ref asn
    3.141000
    3.140000
    <Reference: 4D65822107FCFD52>
    > pi_ref_ref@@
    3.141000
    3.141000
    3.140000
    <Reference: 4D65822107FCFD52>

### Code Blocks

JSL does not have functions, per se. Rather, you can define and execute code blocks. Code blocks are defined by placing code between curly brackets ({ and }).

    > { 3 4 + }
    <CodeBlock>

To execute a code block on the stack, use the `!` operator.

    > { 3 4 + }
    <CodeBlock>
    > !
    7.000000

Of course, you could also write the following equivalent code:

    > {3 4 +}!         
    7.000000

Code blocks can be assigned to variables.

    > { 3 4 + } 'add_three_and_four asn
    > add_three_and_four!
    7.000000

You can use code blocks to simulate functions by passing parameters on the stack.

    > { + } 'add asn
    > 3 4 add!
    7.000000

Of course, you can even pass another code block as a parameter.

    > { ! } 'perform_operation asn 
    > 3 4 { + } perform_operation!
    7.000000

Code blocks can also be nested. Each nested code block has its own variable scope, which inherits from that of its parent.

    > { { 3 4 +}! 2 *}!
    14.000000
    > { 3.14 'pi asn { pi 2 *}! }!
    6.280000
    14.000000
    > { 3.14 'pi asn { pi 2 * 'two_pi asn}! two_pi }!
    Error: Variable 'two_pi' undefined in the local scope.

### Control Flow

JSL has a boolean type that results from comparisons.

    > 3 4 =
    false
    
    > 3 3 =
    true
    
    > 3 4 <
    true

    > 5 6 >
    false

Comparison of values with different types will always result in `false`.

    > "3" false >
    false

Comparisons involving greater than or less than (`<`, `<=`, `>`, and `>=`) will result in an error if both operands are boolean.

    > false true >
    Error: Operation greater than cannot be applied to boolean.

Negation can be achieved using the `~` operator.

    > true~
    false

    > 2 3 =
    false

    > 2 3 =~
    true

You can use `if` to control whether a code block is executed.

    > { 1 } 1 1 = if
    1.000000

    > { 2 } 1 1 =~ if
    (No output)

JSL also supports for loops. The following code computes the sum of the numbers from 1 to 100.

    > 0 'sum asn
    > { 1 'i asn } { i 100 <= } { sum i + 'sum asn } { i 1 + 'i asn } for 
    > sum
    5050.000000

The first code block is the for loop initializer. Here we set the counter variable `i`. The second code block should place a boolean onto the stack (`true` if the body is to be executed, and `false` if not). The third code block is the body. Finally, the fourth code block contains any post-body computations, usually updating the counter variable.

Here is a function that uses for loops to calculate the mean of numbers on the stack:

    {
        'n asn

        { 0 'i asn } { i n 1 - < }
        {
            +
        }
        { i 1 + 'i asn } for

        n /
    } 'mean asn

You can use it as follows:

    > 2 4 2 mean!
    3.000000
    > 10 13 15 3 mean!
    12.666667
    3.000000

We couuld even calculate, for example, the mean of the numbers 1...100:

    > { 1 'i asn } { i 100 <= } { i } { i 1 + 'i asn } for
    100.000000
    99.000000
    98.000000
    97.000000
    96.000000
    95.000000
    94.000000
    93.000000
    92.000000
    ...
    1.000000

    > { 1 'i asn } { i 100 <= } { i } { i 1 + 'i asn } for 100 mean!
    50.500000

### Recursion

JSL supports recursive code blocks. Here is an implementation of factorial:

    {
        'n asn

        {
            1
        } n 1 = if

        {
            n n 1 - fac! *
        } n 1 > if
    } 'fac asn

Tail recursion is also possible.

    {

        {
            'acc asn
            'n asn

            {
                acc
            } n 1 = if

            {
                n 1 -
                acc n *
                fac_tr!
            } n 1 > if
        } 'fac_tr asn

        1 fac_tr!
    } 'fac asn

### Stack Manipulation

You can duplicate the top-most item on the stack with the `dup` operation.

    > 2 dup
    2.000000
    2.000000

You can use this to implement a squaring function:

    > { dup * } 'square asn
    > 16 square!
    256.000000
    > 16 square! square!
    65536.000000
    256.000000

If you want to discard the top-most item on the stack, use the `pop` operation.

    > 3
    3.000000
    > pop
    > 

### Lists

JSL includes an implementation of lists which closely resembles that of functional languages like Lisp or OCaml. First, you start with the empty list:

    > <>
    <>

Then, you can add elements onto the list using the `::`, or cons operator. This adds the object on the top of the stack to a list, and returns a new list. Lists are immutable in JSL.

    > <> 3 ::
    <...>
    > 4 ::
    <...>

To get at elements within the list, we can use the split operation. This returns the first element of the list (the head) on the top of the stack, and a new list containing the remaining elements (the tail) below the head:

    > <> 3 :: 4 :: split
    4.000000
    <...>
    > pop split
    3.000000
    <>

Trying to split the empty list will generate an error:

    > <> split
    Error: Unable to split an empty list.

You can use the `empty?` operator to determine if a list is empty:

    > <> empty?
    true

    > <> 3 :: empty?
    false

Here are a couple common list operations implemented in JSL. The first procedure reverses a list. The second function returns a new list containing the result of a procedure run on each of the input list elements:

    {
        {
            'acc asn
            'list asn

            { acc } list empty? if
            {
                list split
                'head asn
                'tail asn

                tail acc head :: rev_tr!
            } list empty? ~ if
        } 'rev_tr asn

        <> rev_tr
    } 'rev asn

    {
        {
            'acc asn
            'list asn
            'func asn

            { acc rev! } list empty? if
            {

                list split
                'head asn
                'tail asn

                func tail acc head func! :: map_tr!

            } list empty? ~ if
        } 'map_tr asn

        <> map_tr!
    } 'map asn

Here's an example using `map` to double each element within a list:

    > { 2 * } <> 3 :: 4 :: map!
    <...>
    > split
    8.000000
    <...>
    > pop split
    6.000000
    <>



