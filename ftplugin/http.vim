" Vim filetype plugin for HTTP files
" Language: HTTP
" Maintainer: httpx Plugin

if exists("b:did_ftplugin")
  finish
endif
let b:did_ftplugin = 1

" Set comment string
setlocal commentstring=#\ %s

" Enable folding for request blocks
setlocal foldmethod=expr
setlocal foldexpr=GetHTTPFold(v:lnum)

function! GetHTTPFold(lnum)
  let line = getline(a:lnum)
  if line =~# '^###'
    return '>1'
  endif
  return '='
endfunction

" Syntax highlighting
if exists("b:current_syntax")
  finish
endif

" Comments
syntax match httpComment "^#.*$" contains=@Spell
syntax match httpComment "^//.*$" contains=@Spell

" Variables
syntax match httpVariable "@\w\+" nextgroup=httpVariableValue
syntax match httpVariableValue "=.*$" contained

" Request separator
syntax match httpSeparator "^###.*$"

" HTTP Methods
syntax keyword httpMethod GET POST PUT DELETE PATCH HEAD OPTIONS TRACE CONNECT

" URLs
syntax match httpUrl "https\?://[^ ]*"

" Headers
syntax match httpHeader "^[A-Za-z-]\+:" nextgroup=httpHeaderValue
syntax match httpHeaderValue ".*$" contained

" JSON in body
syntax region httpJsonBody start="{" end="}" contains=httpJsonKey,httpJsonString,httpJsonNumber,httpJsonBoolean,httpJsonNull fold
syntax region httpJsonBody start="\[" end="\]" contains=httpJsonKey,httpJsonString,httpJsonNumber,httpJsonBoolean,httpJsonNull fold
syntax match httpJsonKey /"[^"]*":\ze/
syntax region httpJsonString start=/"/ skip=/\\"/ end=/"/
syntax match httpJsonNumber /\<-\?\d\+\(\.\d\+\)\?\([eE][+-]\?\d\+\)\?\>/
syntax keyword httpJsonBoolean true false
syntax keyword httpJsonNull null

" Highlighting
highlight default link httpComment Comment
highlight default link httpVariable Identifier
highlight default link httpVariableValue String
highlight default link httpSeparator Title
highlight default link httpMethod Keyword
highlight default link httpUrl Underlined
highlight default link httpHeader Label
highlight default link httpHeaderValue String
highlight default link httpJsonKey Identifier
highlight default link httpJsonString String
highlight default link httpJsonNumber Number
highlight default link httpJsonBoolean Boolean
highlight default link httpJsonNull Constant

let b:current_syntax = "http"