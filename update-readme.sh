#!/bin/bash
godoc2md -links=false . >README.md
sed -i 's/import "."/import "github.com\/bitfield\/cronrun"/g' README.md
sed -i 's/\/src\/target\///g' README.md
