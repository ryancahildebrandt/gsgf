# GSGF

## *JSpeech Grammar Format, now in Go*

---

[*Open*](https://gitpod.io/#https://github.com/ryancahildebrandt/gsgf) *in gitpod*

## Purpose

This project allows you to generate natural language expressions (productions) from context free grammars in JSGF and JJSGF formats.

---

## Background

### CFG

### \[J\]JSGF

<https://www.infragistics.com/help/winforms/ig-spe-ebnf-file-format-overview>
<https://www.w3.org/TR/jsgf/#14185>
<https://support.voicegain.ai/hc/en-us/articles/360048936511-JJSGF-Grammars>

### Similar Tools

---

## Approach

### Theory

### Algorithm

### Code

### Implementation Notes

---

## Usage

```shell
# show command or general help (-h flag optional)
gsgf [generate|sample|export] [-h]

# run on file with default options
gsgf [generate|sample|export] example.jsgf

```

---

## Outputs

- [gsgf](./gsgf) executable

---

## References

Schlimm, D. (2022). Tables as Powerful Representational Tools. In: Giardino, V., Linker, S., Burns, R., Bellucci, F., Boucheix, JM., Viana, P. (eds) Diagrammatic Representation and Inference. Diagrams 2022. Lecture Notes in Computer Science(), vol 13462. Springer, Cham. <https://doi.org/10.1007/978-3-031-15146-0_15>

---

// go might have widespread use of structs and methods as a standin for classes but i think structs for data + standalone functions is gonna be better, procedural amap
// let data be data don't use the struct as a container for logic

White space is ignored before the definition, between the public keyword and the rulename, around the equal-sign character, and around the semi-colon.
White space is significant within the rule expansion.
| operator applies to the entire group it's in, NOT just the expression to immediate right or left
every edgelist is constructed such that it has exactly one start and end node, so edges in e can be updated using 1 start and 1 end node from a
tabs and newlines need to be double escaped and will be included in productions as \t and \n
"//" can be present in expressions but will not be parsed as a weight or a comment
edges have default weight 1.0
public/private doesnt matter for imports, just for productions
grammars can import from any subdir
grammars in subdirs can import from parent dirs IF they're both in the subdirs of the main namespace
all rules from imported grammars are imported to simplify dependency resolution, specifying a rule is only useful for specificity to the user
allows for import <gram> w no rule or *
a complete grammar will resolve even with invalid imports
*+ etc do not function as they do in regex, have no special meaning here
no d2 rendering built in, to keep it lightweight
dir must exist if provided in outFile
