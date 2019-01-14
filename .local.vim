edit alg.go
"edit toknify/toknify.go
edit expr/expr.go
edit expr/translater.go
edit util/util.go
edit expr/expr_test.go
edit expr/parse_test.go
edit $GOPATH/src/scratch/scratch.go
edit ./docs/design.txt
edit .local.vim
edit ~/.vimrc
buffer 1

nnoremap <leader>r :!clear<cr>:write<cr>:!go run %<cr>
nnoremap <leader>t :!clear<cr>:write<cr>:!go test ./...<cr>

