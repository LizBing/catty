import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;

public class ProxyObjectMethods {
    interface Marker {}

    static class Handler implements InvocationHandler {
        Object self;

        @Override
        public Object invoke(Object proxy, Method method, Object[] args) {
            self = proxy;
            if (method.getName().equals("toString")) {
                return "proxy-text";
            }
            if (method.getName().equals("hashCode")) {
                return Integer.valueOf(12345);
            }
            if (method.getName().equals("equals")) {
                return Boolean.valueOf(proxy == args[0]);
            }
            return null;
        }
    }

    public static void main(String[] args) {
        Handler handler = new Handler();
        Marker marker = (Marker) Proxy.newProxyInstance(
                ProxyObjectMethods.class.getClassLoader(),
                new Class<?>[] {Marker.class}, handler);
        System.out.println(marker.toString());
        System.out.println(marker.hashCode());
        System.out.println(marker.equals(marker));
        System.out.println(marker.equals(new Object()));
        System.out.println(handler.self == marker);
    }
}
