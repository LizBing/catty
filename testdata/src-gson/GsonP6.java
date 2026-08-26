import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import java.util.List;

public class GsonP6 {
    public static void main(String[] args) {
        Gson gson = new Gson();
        try {
            java.util.List<String> list = gson.fromJson("[\"a\",\"b\"]",
                new TypeToken<java.util.List<String>>(){}.getType());
            System.out.println("list_size=" + list.size());
        } catch (Throwable t) {
            System.out.println("typeToken=" + t.getClass().getSimpleName()
                + ":" + t.getMessage().substring(0, Math.min(60, t.getMessage().length())));
        }
    }
}
