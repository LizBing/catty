package catty.runtime;

/** Stable admission failure raised by Catty. */
public final class PluginAdmissionException extends RuntimeException {
    private final String code;

    public PluginAdmissionException(String code) {
        super(code);
        this.code = code;
    }

    public String getCode() {
        return code;
    }
}
