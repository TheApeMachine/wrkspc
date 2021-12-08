call plug#begin('~/.vim/plugged')
Plug 'ervandew/supertab'
Plug 'sheerun/vim-polyglot'
Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
Plug 'SirVer/ultisnips'
Plug 'itchyny/lightline.vim'
Plug 'dyng/ctrlsf.vim'
Plug 'kien/ctrlp.vim'
Plug 'mhinz/vim-signify'
Plug 'morhetz/gruvbox'
Plug 'preservim/vimux'
Plug 'tpope/vim-commentary'
Plug 'voldikss/vim-floaterm'
Plug 'embark-theme/vim', { 'as': 'embark' }
Plug 'honza/vim-snippets'
Plug 'tpope/vim-fugitive'
Plug 'dense-analysis/ale'
Plug 'maximbaz/lightline-ale'
Plug 'justinmk/vim-sneak'
Plug 'dracula/vim', { 'as': 'dracula' }
call plug#end()

syntax on

set modifiable
set belloff=all
set visualbell
set encoding=utf-8
set fileencoding=utf-8
set number relativenumber
filetype plugin indent on
set backspace=indent,eol,start
set cmdheight=2
set updatetime=750
set shortmess+=c
set undodir=~/.vim/undodir

if has("nvim-0.5.0") || has("patch-8.1.1564")
  set signcolumn=number
else
  set signcolumn=yes
endif

function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1]  =~# '\s'
endfunction

set rtp+=/usr/local/opt/fzf
let g:go_list_type = "quickfix"

set clipboard=unnamed
let mapleader=" "
nnoremap ; :
nnoremap , n
inoremap jk <ESC>

nnoremap rt :FloatermSend go run main.go run<cr>
nnoremap re :FloatermSend go mod tidy<cr>

nnoremap dlv :GoDebugStart .<cr>
let g:floaterm_keymap_new    = '<F7>'
let g:floaterm_keymap_prev   = '<F8>'
let g:floaterm_keymap_next   = '<F9>'
let g:floaterm_keymap_toggle = '<F10>'
let g:floaterm_keymap_kill   = '<F12>'

nnoremap spc :setlocal spell spelllang=en_us<cr>
nnoremap fsc :set nospell<cr>

nnoremap gb :Git blame<cr>
nnoremap gh :Git diff<cr>
nnoremap hh :Git log<cr>
nnoremap gt :Git<cr>
nnoremap gp :Git push<cr>

nnoremap <silent>kl :vertical resize +5<CR>
nnoremap <silent>lk :vertical resize -5<CR>
nnoremap <silent>hj :resize +5<CR>
nnoremap <silent>jh :resize -5<CR>

let g:ctrlp_user_command = ['.git/', 'git --git-dir=%s/.git ls-files -oc --exclude-standard']

if executable('rg')
  set grepprg=rg\ --color=never
  let g:ctrlp_user_command = 'rg %s --files --color=never --glob ""'
  let g:ctrlp_use_caching = 0
else
  let g:ctrlp_clear_cache_on_exit = 0
endif

nnoremap <leader>/ :vsplit<cr>
nnoremap <leader>. :split<cr>
nnoremap <leader>, :Ex<cr>
nnoremap <leader>o :only<cr>
nnoremap <leader>m <C-w><C-r><c-w>h
nnoremap <c-j> <c-w>j
nnoremap <c-k> <c-w>k
nnoremap <c-h> <c-w>h
nnoremap <c-l> <c-w>l
nnoremap <leader>h <c-w>K
nnoremap <leader>v <c-w>H
nmap <leader>zx <Plug>CtrlSFPrompt
map K :GoDoc<cr>
map <C-n> :cnext<cr>
map <C-m> :cprevious<cr>
nnoremap <leader>a :cclose<cr>
nnoremap gml :GoMetaLinter<cr>
nnoremap gcl :GoCallees<cr>
nnoremap dfs :GoDefStack<cr>
nnoremap rnm :GoRename
nnoremap vrv :'<,'>GoFreevars<cr>
nnoremap nsp :e /home/theapemachine/.vim/plugged/vim-snippets/UltiSnips/go.snippets<cr>
nnoremap vmr :VimuxRunCommand<cr>
map vmt :wa<cr> :GolangTestCurrentPackage<cr>
nnoremap <leader>q :quit<CR>

let g:multi_cursor_use_default_mapping = 0
let g:multi_cursor_start_word_key      = '<A-s>'
let g:multi_cursor_select_all_word_key = '<A-a>'
let g:multi_cursor_next_key            = '<A-n>'
let g:multi_cursor_prev_key            = '<A-p>'
let g:multi_cursor_skip_key            = '<A-x>'
let g:multi_cursor_quit_key            = '<Esc>'

let g:go_metalinter_autosave = 1

func! s:lines_count()
  echom line('$') . ' lines in buffer'
endfunc

func! s:cmdMenu()
  call popup_menu([
        \ 'rt     (run)              | re  (tidy)          | dif (delete in func)',
        \ 'hlp    (this menu)        | gb  (git blame)     | dfs (def stack)',
        \ 'nsp    (new snippet)      | gh  (git diff)      | ]]  (next func)',
        \ 'fk     (easy motion)      | hh  (git log)       | [[  (prev func)',
        \ 'kl     (vertical resize+) | F7  (terminal)      | rnm (go rename)',
        \ 'lk     (vertical resize-) | F12 (kill terminal) | C-n (next error)',
        \ '<spc>m (rotate splits)    | dlv (go debug)      | C-m (prev error)',
        \ 'gml    (go meta linter)   | gcl (go callees)    | vrv (extract func)',
        \ 'rt     (go run)           | re  (go mod tidy)   | spc (spellcheck)',
        \ 'C-j    (ale next)         | C-k (ale previous)  | fsc (fuck spellcheck)',
        \], #{callback: ''})
endfunc

nnoremap hlp :call <SID>cmdMenu()<cr>

set completeopt=menu,longest
let g:SuperTabDefaultCompletionType = "<c-n>"
let g:UltiSnipsExpandTrigger="<tab>"
let g:UltiSnipsJumpForwardTrigger="<tab>"
let g:UltiSnipsJumpBackwardTrigger="<s-tab>"

autocmd BufNewFile,BufRead *.go setlocal noexpandtab tabstop=4 shiftwidth=4

let g:ale_sign_error = '💔'
let g:ale_sign_warning = '👀'
let g:ale_fixers = {
\   '*': ['remove_trailing_lines', 'trim_whitespace'],
\   'javascript': ['prettier'],
\}
let g:ale_fix_on_save = 1
let g:ale_completion_enabled = 1
let g:ale_completion_autoimport = 1
let g:ale_set_loclist = 0
let g:ale_set_quickfix = 1
let g:ale_lint_on_text_changed = 'never'
let g:ale_lint_on_insert_leave = 0
let g:ale_lint_on_enter = 0

nmap <silent> <C-k> <Plug>(ale_previous_wrap)
nmap <silent> <C-j> <Plug>(ale_next_wrap)

autocmd InsertLeave * write
if version >= 702
  autocmd BufWinLeave * call clearmatches()
endif

let g:go_highlight_extra_types=1
let g:go_highlight_operators=1
let g:go_highlight_functions=1
let g:go_highlight_function_parameters=1
let g:go_highlight_function_calls=1
let g:go_highlight_types=1
let g:go_highlight_fields=1
" let g:go_highlight_format_strings=1
" let g:go_highlight_variable_declarations=1
" let g:go_highlight_variable_assignments=1

let &t_8f="\<Esc>[38;2;%lu;%lu;%lum"
let &t_8b="\<Esc>[48;2;%lu;%lu;%lum"
set termguicolors

set t_Co=256
set background=dark
colorscheme gruvbox

aug QFClose
  au!
  au WinEnter * if winnr('$') == 1 && &buftype == "quickfix"|q|endif
aug END
