package shell

const zshHookTemplate = `
__xc_session_file() {
    echo "/tmp/xc_${USER}_$$.tmp"
}

__xc_cmd=""
__xc_outfile=""
__xc_capturing=0

__xc_preexec() {
    local cmd="$1"
    [[ "$cmd" == xc* ]] && return
    [[ "$cmd" == __xc_* ]] && return

    __xc_cmd="$cmd"
    __xc_outfile="$(mktemp /tmp/xc_out_XXXXXX)"
    __xc_capturing=1

    exec 3>&1 4>&2
    exec 1> >(tee -a "$__xc_outfile") 2>&1
}

__xc_precmd() {
    if [[ "$__xc_capturing" -eq 1 ]]; then
        __xc_capturing=0
        exec 1>&3 2>&4
        exec 3>&- 4>&-

        local session_file
        session_file="$(__xc_session_file)"
        {
            print "CMD:${__xc_cmd}"
            print "OUTPUT:"
            if [[ -f "$__xc_outfile" ]]; then
                cat "$__xc_outfile"
                rm -f "$__xc_outfile"
            fi
        } > "$session_file"
        chmod 600 "$session_file"
    fi
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec __xc_preexec
add-zsh-hook precmd __xc_precmd
`

func ZshHook() string {
	return zshHookTemplate
}
