# Relational Data Model
This library defines a simple parent-child relational database with more references supported.

Start with a json file to describe your model, then generate go code and write the missing parts yourself...

# Install
```
go install github.com/jansemmelink/model/modelgen
```

# Usage
```
% modelgen -f mymodel.json -o ./mymodel
```
That will generate code in ./mymodel/api and ./mymodel.system
