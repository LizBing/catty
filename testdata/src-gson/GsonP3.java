import com.google.gson.Gson;

public class GsonP3 {
    public static class Pojo {
        public String name;
        public int count;
    }

    public static void main(String[] args) {
        Gson gson = new Gson();
        Pojo p = gson.fromJson("{\"name\":\"hello\",\"count\":42}", Pojo.class);
        System.out.println("name=" + p.name + " count=" + p.count);
    }
}
