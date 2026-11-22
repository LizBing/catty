package catty.runtime;

public final class PluginAdmissionException extends RuntimeException {
    private final String code;

    public PluginAdmissionException(String code) {
        super(code);
        this.code = code;
    }

    public PluginAdmissionException(String code, String message) {
        super(message);
        this.code = code;
    }

    public String getCode() {
        return code;
    }
}
