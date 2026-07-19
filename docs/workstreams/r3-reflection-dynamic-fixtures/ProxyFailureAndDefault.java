import java.io.IOException;
import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.lang.reflect.UndeclaredThrowableException;

public class ProxyFailureAndDefault {
    interface Service {
        default String defaultName() {
            return "default";
        }

        void unchecked();
        void checked();
    }

    static class Handler implements InvocationHandler {
        @Override
        public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
            if (method.isDefault()) {
                return InvocationHandler.invokeDefault(proxy, method, args);
            }
            if (method.getName().equals("unchecked")) {
                throw new IllegalStateException("unchecked");
            }
            throw new IOException("checked");
        }
    }

    public static void main(String[] args) {
        Service service = (Service) Proxy.newProxyInstance(
                ProxyFailureAndDefault.class.getClassLoader(),
                new Class<?>[] {Service.class}, new Handler());
        System.out.println(service.defaultName());
        try {
            service.unchecked();
        } catch (IllegalStateException expected) {
            System.out.println(expected.getMessage());
        }
        try {
            service.checked();
        } catch (UndeclaredThrowableException expected) {
            System.out.println(expected.getUndeclaredThrowable().getClass().getName());
            System.out.println(expected.getUndeclaredThrowable().getMessage());
        }
    }
}
