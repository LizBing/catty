import com.google.gson.Gson;

public class GsonP4 {
    public static class Address {
        public String city = "unknown";
        public String zip = "00000";
    }

    public static class Person {
        public String name = "test";
        public Address addr = new Address();
        public int age = 30;
    }

    public static void main(String[] args) {
        Gson gson = new Gson();
        Person p = new Person();
        p.name = "alice";
        p.addr.city = "tokyo";
        String json = gson.toJson(p);
        System.out.println("json=" + json);
    }
}
