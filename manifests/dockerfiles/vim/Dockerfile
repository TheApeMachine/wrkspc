FROM theapemachine/zsh:v1.0

# Since this uses the zsh image as a base, we're already dropped to a normal user, keep
# that in mind and use sudo where needed.
RUN sudo apk add --no-cache ca-certificates vim bind-tools nodejs curl

USER $USER
WORKDIR $HOME

# Install a vim plugin manager.
RUN curl -fLo ~/.vim/autoload/plug.vim --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
# Copy any .vimrc or other config files present in the context.
COPY . .
# Run the installers for all plugins.
# RUN vim +PlugInstall +qall
RUN vim -V -es -u ~/.vimrc -i NONE -c "PlugInstall" -c "qa"

# Start working...
ENTRYPOINT /bin/zsh
