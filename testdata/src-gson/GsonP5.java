import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;
import java.lang.annotation.Annotation;
import java.lang.reflect.Field;

public class GsonP5 {
    public static class Renamed {
        @SerializedName("user_name")
        public String userName = "test";

        @SerializedName("age_years")
        public int age = 25;
    }

    public static void main(String[] args) throws Exception {
        Gson gson = new Gson();
        // Without annotation support, gson falls back to field name
        String json = gson.toJson(new Renamed());
        System.out.println("json=" + json);

        // Probe: does getAnnotation return non-null?
        Field f = Renamed.class.getDeclaredField("userName");
        Annotation a = f.getAnnotation(SerializedName.class);
        System.out.println("annotation=" + (a != null ? "found" : "null"));
    }
}
