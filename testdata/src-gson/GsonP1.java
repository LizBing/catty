import com.google.gson.Gson;

public class GsonP1 {
    public static void main(String[] args) throws Exception {
        Class<?> c = Class.forName("com.google.gson.Gson");
        System.out.println("class=" + c.getName());
        Object g = c.getDeclaredConstructors()[0].newInstance();
        System.out.println("gson=" + (g != null));
    }
}
