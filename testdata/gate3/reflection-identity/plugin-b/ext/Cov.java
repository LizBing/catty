package ext;

import app.Factory;

public class Cov extends Factory {
    public static Cov create() { return new Cov(); }

    @Override
    public String make() { return "abc"; }
}
