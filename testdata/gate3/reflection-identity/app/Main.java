package app;

import catty.runtime.Plugins;
import java.lang.reflect.Method;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader l1 = Plugins.register(Plugins.readJar(args[0]), app);
        ClassLoader l2 = Plugins.register(Plugins.readJar(args[1]), app);
        Class<?> w1 = Class.forName("ext.Worker", false, l1);
        Class<?> w2 = Class.forName("ext.Worker", false, l2);
        Class<?> t1 = Class.forName("ext.Token", false, l1);
        Class<?> t2 = Class.forName("ext.Token", false, l2);
        System.out.println(t1 == t2);

        Method newToken = t1.getMethod("create", new Class[0]);
        Object token1 = newToken.invoke(null);
        Method newOther = t2.getMethod("create", new Class[0]);
        Object token2 = newOther.invoke(null);
        Method accept = w1.getMethod("accept", new Class[]{t1});
        System.out.println(((Integer) accept.invoke(null, token1)).intValue());
        try {
            w1.getMethod("accept", new Class[]{t2});
            System.out.println(0);
        } catch (NoSuchMethodException e) {
            System.out.println(1);
        }
        try {
            accept.invoke(null, token2);
            System.out.println(0);
        } catch (IllegalArgumentException e) {
            System.out.println(2);
        }

        Class<?> a1 = Class.forName("[Lext.Token;", false, l1);
        Class<?> a2 = Class.forName("[Lext.Token;", false, l2);
        Method array = w1.getMethod("accept", new Class[]{a1});
        Object tokens1 = w1.getMethod("tokens", new Class[0]).invoke(null);
        Object tokens2 = w2.getMethod("tokens", new Class[0]).invoke(null);
        System.out.println(((Integer) array.invoke(null, tokens1)).intValue());
        try {
            array.invoke(null, tokens2);
            System.out.println(0);
        } catch (IllegalArgumentException e) {
            System.out.println(3);
        }

        Class<?> cov = Class.forName("ext.Cov", false, l1);
        Object covObject = cov.getMethod("create", new Class[0]).invoke(null);
        Method make = cov.getMethod("make", new Class[0]);
        Object made = make.invoke(covObject);
        System.out.println(((String) made).length());
    }
}
