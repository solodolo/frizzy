# frizzy

frizzy is a(nother) static site generator written in Go. 

It uses a custom templating language that is parsed using a totally overkill LR1 parser.

## Templates
todo: finish this section

### for loops

```
  {{for foo in "bar"}}
    <h1>{{: foo.title}}<h1>
  {{end}}
```

### if statements
```
  {{if a < b}}
    <p>a is smaller</p>
  {{else_if a > b}}
    <p>b is smaller</p>
  {{else}}
    <p>a equals b</p>
  {{end}}
```

### variable assignment
```
  {{title = "this is the title"}}
  .
  .
  .
  {{: title}}
```

### function calls
```
  {{ paginator()}}
```

## Installation
todo

## Usage
todo
### Configuration
todo

## Contributing
todo