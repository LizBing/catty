import com.google.gson.Gson;

public class GsonP2 {
    public static class Pojo {
        public String name = "catty";
        public int count = 7;
        public boolean flag = true;
    }

    public static void main(String[] args) {
        Gson gson = new Gson();
        String json = gson.toJson(new Pojo());
        System.out.println(json);
    }
}
