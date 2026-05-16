package shell

const bashHookTemplate = `
__xc_session_file() {
    local dir="${XDG_RUNTIME_DIR:-/tmp}"
    echo "${dir}/xc_$(id -u)_$$.tmp"
}

__xc_cmd=""
__xc_capturing=0
__xc_outfile=""

__xc_preexec() {
    local cmd="$__xc_last_cmd"
    [[ -z "$cmd" ]] && return
    [[ "$cmd" == __xc_* ]] && return
    [[ "$cmd" == xc* ]] && return

    __xc_cmd="$cmd"
    __xc_outfile="$(mktemp "${XDG_RUNTIME_DIR:-/tmp}/xc_out_XXXXXX")"
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
            echo "CMD:${__xc_cmd}"
            echo "OUTPUT:"
            if [[ -f "$__xc_outfile" ]]; then
                cat "$__xc_outfile"
                rm -f "$__xc_outfile"
            fi
        } > "$session_file"
        chmod 600 "$session_file"
    fi
    __xc_last_cmd=""
}

__xc_debug_trap() {
    if [[ "$BASH_COMMAND" != __xc_* ]] && [[ -z "$__xc_last_cmd" ]]; then
        # Use history for the full pipeline text; fall back to BASH_COMMAND for the first word
        __xc_last_cmd="$(HISTTIMEFORMAT= history 1 | sed 's/^[[:space:]]*[0-9]*[[:space:]]*//')"
        [[ -z "$__xc_last_cmd" ]] && __xc_last_cmd="$BASH_COMMAND"

        # Write command immediately so xc can read it when running inside a pipeline
        local session_file
        session_file="$(__xc_session_file)"
        { echo "CMD:${__xc_last_cmd}"; echo "OUTPUT:"; } > "$session_file"
        chmod 600 "$session_file"
    fi
}

trap '__xc_debug_trap' DEBUG
PROMPT_COMMAND="__xc_precmd${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
`

func BashHook() string {
	return bashHookTemplate
}
