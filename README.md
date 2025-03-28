# GSGF

## *JSpeech Grammar Format, now in Go*

---

[*Open*](https://gitpod.io/#https://github.com/ryancahildebrandt/gsgf) *in gitpod*

## Purpose

This project provides utilities for working with JSGF and JJSGF context free grammars, including producing natural language expressions and parsing grammar files into graph representations.

---

## Background

Building something like this from scratch was admittedly outside my wheelhouse, and it required picking up a lot of techniques from other fields. Even with the multitude of learning resources available on the internet, this was the most challenging project I've done and I learned quite a lot by doing it. Throughout the process, I found no shortage of helpful resources explaining individual pieces needed for this project (especially from Kay Lack [@neoeno](https://github.com/neoeno) and her [Youtube Channel](https://www.youtube.com/neoeno4242)), but putting those pieces together was a huge challenge. I've collected the most helpful information to provide an overview of the moving parts of the project and how they fit together.

At a high level, the program:

1) Reads in a grammar file, storing rules and imports
2) Splits each rule into an array of individual tokens
3) Iterates over the token array to construct finite state automaton in the form of a directed graph
4) Calculates all paths between the initial and final graph nodes
5) Maps each traversal path to the original array of tokens to produce natural language expressions

If 100% of that makes sense to you, congrats! You probably have very little to learn from reading the rest of this.
If not, check the sections below for more details.

### Finite State Automata [<sup>1, 2, 3</sup>](#references)

As the name suggests, a finite state automaton (FSA) is an abstract system which can exist in a limted number of states. States can be changed via some input or process over time, and a change from one state to another is called a transition. In practice, specific states are often designated as the initial/final states, so that we can keep track of what state we're starting from and what state corresponds to some success condition. Regular expressions, computational states, and (most relevant here) context free grammars can all be represented as FSA.

### Context Free Grammars [<sup>4, 5, 6</sup>](#references)

Context free grammars define the symbols that can appear in a given "language".  A language in this case can refer to natural languages, programming languages, or really any sequence of characters with definable patterns. Grammars are made up of rules, consisting of an identifier (often called a nonterminal symbol) and an expansion, which consists of 0 or more identifiers and/or string literals (often called a terminal symbols).

There are several common formats for grammar files, including Backus-Naur Form, Extended Backus-Naur Form, and JSpeech Grammar Format. Different formats may have symbols that allow for more complex relationships between terminal symbols (grouping, alternates, repetition, etc.), but generally speaking any context free grammar format can represent any context free language.

### \[J\] JSGF [<sup>7, 8</sup>](#references)

JSGF (and its json wrapped cousin JJSGF) is a grammar representation format used in Java based systems. While commonly used in speech **recognition** systems, the same underlying grammar representation can be used to **produce** valid expressions.

The JSGF format draws some syntactic conventions from the Java programming language, including valid characters and identifiers. JSGF grammars also inherit Java's module/package import and export mechanics, meaning that rules can both be 1) imported from other grammars and/or 2) designated as public/private (visible or invisible to other grammars). This import system allows for organizing rules into meaningful groups in different files, and for multiple grammars to import the same rule from an external grammar.

---

## Approach

Below is some pseudocode for how the program works, at a slightly lower level than the bullet list at the top of the page.

<pre><code>
Load <b>grammar</b> from file with FromJSGF(<b>file</b>)
 - For each <b>line</b> in <b>file</b>:
  - Create <b>rule</b> with ParseRule(<b>line</b>)
  - Split <b>rule</b> into <b>tokens</b> with ToTokens(<b>rule</b>)
  - Convert <b>tokens</b> to an edgelist/graph with ToEdgeList(<b>tokens</b>)
 - If <b>grammar</b> is not complete:
  - Collect imports from <b>grammar</b> with CreateNameSpace
  - Add rules to <b>grammar</b> with ImportNameSpace
 - Combine rules in order with ResolveRules(<b>grammar</b>)
  - For each <b>rule</b> in <b>grammar</b>, insert referenced rules with ResolveReferences(<b>rule</b>):
   - For each <b>reference</b> in <b>rule</b>, insert referenced <b>rule</b> with singleResolveReference(<b>rule</b>, <b>reference</b>)
    - For each occurence of <b>reference</b> in target <b>rule</b>, compose <b>ref graph</b> into <b>rule graph</b> with composeGraphs(<b>rule graph</b>, <b>ref graph</b>)
 - Collect all productions from the resolved <b>grammar</b> with GetAllProductions(<b>grammar</b>)
  - For each public <b>rule</b> in<b>*grammar</b>:
   - Collect productions with getProductions(<b>rule</b>)
    - Traverse <b>rule graph</b> with getAllPAths(<b>graph</b>)
    - Remove tokens that do not contribute to productions with filterTokens(<b>tokens</b>)
    - For each <b>path</b> in <b>graph</b>:
     - Map <b>path</b> to filtered <b>tokens</b> with getSingleProduction(<b>path</b>, **tokens</b>)
</code></pre>

### My Trouble With Parsing [<sup>9, 10</sup>](#references)

The main thing missing from the resources I was albe to find (and my non-CS background) was how to implement a lot of these above pieces in code. Pretty early on I hit a roadblock with my understanding of going from a tokenized expression to a useful representation in code. Most explanations focused on recursive descent parsing to construct a parse tree, and then using operations on tokens within the expression to produce some output. There's no shortage of resources on recursive descent parsing, so going from an array of tokens to a parse tree made enough sense. Even the FSA bit was pretty straightforward conceptually, making an array of tokens into a graph containing all possible state transitions. The part that didn't really make sense to me was where to go from there.

My understanding of grammars/FSA/expression production really hinged on the graph data structure, and it wasn't apparent how to go from tokens to a graph, or a parse tree to a graph, or skip the graph altogether and get prodductions directly from the parse tree. It made more sense (given my understanding of the problem) to go straight from an array of tokens to a graph representation if possible.

Most of what I could find either 1) stopped at evaluating simple math expressions or 2) was written for people already well versed in compiler design, which I am very much not. Either way I didn't have a reference for how to approach the problem, so I decided on the approach that made the most sense to me and that I felt I could at least make a decent go at without losing my mind.

And it happened to work! I'd say the solution I came up with was probably most similar to the McNaughton-Yamada-Thompson algorithm, though the only explanations of that algorithm I could find were in pseudocode and described a recursive approach.

Thankfully I was able to get a working iterative algorithm that:

1) Accepts an array of tokens
2) Iterates over that array exactly once
3) Doesn't rely on special classes/structs, instead using 2 int pointers, 1 pseudo-stack, and 1 map to keep track of valid state transitions
4) Returns an array of edges representing all valid transitions in the FSA

I have no reason to believe my approach is in any way novel, but I'm still proud of it and will toot my own horn at least a little via the following detailed walkthrough ~ ~ with pictures ~ ~.

### Algorithm

```go
// Converts a slice of tokens/Expressions to an EdgeList
// Uses flow control tokens (), [], | to capture possible state transitions between tokens
// Every edgelist is constructed such that it has exactly one start and end node
func ToEdgeList(arr []Expression) EdgeList {
 var (
  edges      EdgeList
  from       int
  group      int
  groupStack []int
  groupMap   = make(map[int][]int)
 )

 for i, token := range arr {
  switch token {
  case "<EOS>":
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
  case "<SOS>":
   from = i
   groupStack = append(groupStack, i)
   groupMap[i] = []int{}
  case ";":
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
   for _, v := range groupMap {
    for _, v1 := range v {
     edges = append(edges, Edge{From: v1, To: i, Weight: 1.0})
    }
   }
   from = i
  case "(", "[":
   groupStack = append(groupStack, i)
   groupMap[i] = []int{}
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
   from = i
  case ")":
   group = groupStack[len(groupStack)-1]
   for _, v := range groupMap[group] {
    edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
   }
   groupStack = groupStack[:len(groupStack)-1]
   delete(groupMap, group)
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
   from = i
  case "]":
   group = groupStack[len(groupStack)-1]
   for _, v := range groupMap[group] {
    edges = append(edges, Edge{From: v, To: i, Weight: 1.0})
   }
   groupStack = groupStack[:len(groupStack)-1]
   delete(groupMap, group)
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
   edges = append(edges, Edge{From: group, To: i, Weight: 1.0})
   from = i
  case "|":
   group = groupStack[len(groupStack)-1]
   groupMap[group] = append(groupMap[group], from)
   from = group
  default:
   edges = append(edges, Edge{From: from, To: i, Weight: 1.0})
   from = i
  }
 }

 return Unique(edges)
}
```

The promised pictures, walking through each step of the algorithm on simple examples. As noted above, valid transitions are tracked in 2 integer pointers (*from* to track the current source node, *group* to track the node corresponding to the current group's opening token), 1 pseudo-stack (*groupStack* to track any open groups), and 1 map (*groupMap* to track of possible alternates within a group).

1) Example string: abc\
![Simple transitions](./data/figures/fig1.svg)

2) Example string: a|b|c\
![Alternates](./data/figures/fig2.svg)

3) Example string: (a|b|c)\
![Grouping](./data/figures/fig3.svg)

4) Example string: a[b]c\
![Optional](./data/figures/fig4.svg)

### Implementation Notes

The parsers used here to read JSGF files differ in a few ways from the original JSGF implementation. Some of these won't change how you use the grammars, some might, and some of them are features of JSGF as used for recognition systems that didn't make a ton of sense for production generation.

- White space is ignored around the = and ; tokens
- \| operator applies to the entire group it's located in, not just the expression to immediate right or left
- Tabs and newlines need to be double escaped and in order to be included in productions as \t and \n
- "//" can be present in expressions but will not be parsed as a weight or a comment
- The "public" declaration before a rule identifier doesnt matter for imports, just for productions. A rule can be imported even if it isn't declared as public
- Grammars can import from any subdirectory
- In the below example directory:
  - The namespace will not be resolvable if gsgf is called on b, because b imports from a grammar in a parent directory
  - The namespace will be resolvable if gsgf is called on main.jsgf or a.jsgf, *even though b imports from a grammar in a parent directory*, because the namespace resolution process does not depend on imported grammars' relationship to each other, only on their relationship to the main grammar

```
main.jsgf
↳ subdir
 ↳ a.jsgf (imports b)
   subdir
    ↳ b.jsgf (imports a)
```

- All rules from imported grammars are read into the namespace to simplify dependency resolution, specifying a rule in an import statement is only useful for specificity/clarity for the user
- Namespace resolution relies on hashmaps, so namespace collisions are very possible and will overwrite rules defined in parent grammars
- It is also possible to import <gram> without specifying a rule or *
- The namespace resolution process checks for grammar completeness (whether a grammar can be fully resolved only using rules defined in the grammar), so a complete grammar will resolve even with invalid import statements
- JSGF includes * and + as quantifiers for recurring tokens, and these functions are not supported here

### Similar Tools

- [jsgf-gen](https://github.com/synesthesiam/jsgf-gen)
- [IntXeger](https://github.com/k15z/IntXeger)
- [xeger](https://github.com/agarciadom/xeger)

---

## Usage

```shell
# show general or command specific help (-h flag optional)
gsgf [generate|sample|export] [-h]

# generate all productions, shuffling the order and writing to myfile.txt
gsgf generate --shuffle --outFile "myfile.txt" example.jsgf

# sample 100 productions, removing initial and terminal spaces and printing to stdout
gsgf sample --nProductions 100 --removeEndSpaces example.jsgf

# export grammar and minimized graph representations to ./myDir/
gsgf export --exportDir "myDir" --minimize example.jsgf

```

---

## Outputs

- [gsgf](./gsgf) executable
- [Example](./data) JSGF and JJSGF grammars

---

## References

- [1] <https://en.wikipedia.org/wiki/Finite-state_machine>
- [2] <https://www.0de5.net/stimuli/regular-expressions-as-finite-automata>
- [3] <https://www.0de5.net/stimuli/fun-and-games-generating-dfas-from-regular-expressions>
- [4] <https://en.wikipedia.org/wiki/Context-free_grammar>
- [5] <https://en.wikipedia.org/wiki/Backus–Naur_form>
- [6] <https://www.cs.rochester.edu/u/nelson/courses/csc_173/grammars/>
- [7] <https://www.w3.org/TR/jsgf/#14185>
- [8] <https://support.voicegain.ai/hc/en-us/articles/360048936511-JJSGF-Grammars>
- [9] <https://www.0de5.net/stimuli/grammars-parsing-and-recursive-descent>
- [10] <https://en.wikipedia.org/wiki/Thompson%27s_construction>
