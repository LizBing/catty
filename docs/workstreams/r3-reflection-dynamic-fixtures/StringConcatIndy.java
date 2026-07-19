public class StringConcatIndy {
    static String combine(String text, int number, long wide, Object value) {
        return text + ":" + number + ":" + wide + ":" + value;
    }

    public static void main(String[] args) {
        System.out.println(combine("catty", 25, 3000000000L, null));
        System.out.println(combine("\uD83D\uDE3A", -1, 0L, "done"));
    }
}
