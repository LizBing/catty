import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;

public class ProxyDispatch {
    interface Service {
        int add(int left, int right);
        String name();
    }

    static class Handler implements InvocationHandler {
        int calls;

        @Override
        public Object invoke(Object proxy, Method method, Object[] args) {
            calls++;
            if (method.getName().equals("add")) {
                return Integer.valueOf(((Integer) args[0]).intValue()
                        + ((Integer) args[1]).intValue());
            }
            return "catty";
        }
    }

    public static void main(String[] args) {
        Handler handler = new Handler();
        Service service = (Service) Proxy.newProxyInstance(
                ProxyDispatch.class.getClassLoader(),
                new Class<?>[] {Service.class}, handler);
        System.out.println(service.add(4, 5));
        System.out.println(service.name());
        System.out.println(handler.calls);
        System.out.println(Proxy.isProxyClass(service.getClass()));
        System.out.println(Proxy.getInvocationHandler(service) == handler);
    }
}
