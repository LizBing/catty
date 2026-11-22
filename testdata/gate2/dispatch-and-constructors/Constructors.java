public final class Constructors {
    static int note(int value) {
        System.out.println(value);
        return value;
    }

    interface Greeter {
        default int speak() {
            return 51;
        }
    }

    interface Loud extends Greeter {
        default int speak() {
            return Greeter.super.speak() + 1;
        }
    }

    static class Quiet implements Greeter {
    }

    static class Shout implements Loud {
    }

    static final class Custom implements Loud {
        public int speak() {
            return 53;
        }
    }

    static final class Both implements Greeter, Loud {
    }

    static final class Direct implements Greeter {
        public int speak() {
            return Greeter.super.speak() + 1;
        }
    }

    static class Animal {
        int sound() {
            return 61;
        }

        int legs() {
            return 4;
        }
    }

    static class Cat extends Animal {
        int sound() {
            return 62;
        }
    }

    static class Dog extends Animal {
        int sound() {
            return 63;
        }
    }

    static final class Kennel {
        Animal resident = new Dog();
    }

    abstract static class Shape {
        abstract int area();

        int scaled() {
            return area() * 2;
        }
    }

    static final class Square extends Shape {
        private final int side;

        Square(int side) {
            this.side = side;
        }

        int area() {
            return side * side;
        }
    }

    static class Root {
        int rootField = note(10);

        {
            System.out.println(11);
        }

        Root() {
            System.out.println(12);
        }

        Root(int seed) {
            this();
            this.rootField = seed;
            System.out.println(13);
        }

        int plain() {
            return 21;
        }

        int hook() {
            return 0;
        }
    }

    static class Base extends Root {
        int baseField;
        int hidden = 11;

        Base(int seed) {
            super(seed);
            System.out.println(22);
            hook();
            this.baseField = 30;
            System.out.println(23);
        }

        int baseHidden() {
            return hidden;
        }
    }

    static final class Derived extends Base {
        int value = note(40);

        {
            System.out.println(41);
        }

        int hidden = 22;
        int other = note(42);

        Derived() {
            super(9);
            System.out.println(43);
        }

        int hook() {
            System.out.println(value == 0 ? 44 : 45);
            return value;
        }

        int derivedHidden() {
            return hidden;
        }
    }

    static final class Holder {
        int value;
    }

    static final class Ticket {
        private final int number;

        private Ticket(int number) {
            this.number = number;
        }

        static Ticket create(int number) {
            return new Ticket(number);
        }

        int number() {
            return number;
        }
    }

    public static final class PublicBox {
        private final int payload;

        public PublicBox(int payload) {
            this.payload = payload;
        }

        public int payload() {
            return payload;
        }
    }

    static void dispatch() {
        System.out.println(100);
        Animal animal = new Cat();
        System.out.println(animal.sound());
        System.out.println(animal.legs());
        Animal[] zoo = new Animal[2];
        zoo[0] = new Cat();
        zoo[1] = new Dog();
        System.out.println(zoo[0].sound());
        System.out.println(zoo[1].sound());
        Kennel kennel = new Kennel();
        System.out.println(kennel.resident.sound());
        Animal named = new Dog();
        System.out.println(named.sound());
        Animal plain = new Animal();
        System.out.println(plain.sound());
    }

    static void defaults() {
        System.out.println(200);
        Greeter quiet = new Quiet();
        System.out.println(quiet.speak());
        Greeter shout = new Shout();
        System.out.println(shout.speak());
        Loud loud = new Shout();
        System.out.println(loud.speak());
        Greeter custom = new Custom();
        System.out.println(custom.speak());
        Both both = new Both();
        System.out.println(both.speak());
        Greeter bothGreeter = new Both();
        System.out.println(bothGreeter.speak());
        Direct direct = new Direct();
        System.out.println(direct.speak());
        Greeter greeter = new Direct();
        System.out.println(greeter.speak());
    }

    static void abstractShapes() {
        System.out.println(300);
        Shape shape = new Square(3);
        System.out.println(shape.area());
        System.out.println(shape.scaled());
        Square square = new Square(4);
        System.out.println(square.area());
        Shape other = new Square(5);
        System.out.println(other.scaled());
    }

    static Derived constructors() {
        System.out.println(400);
        Derived derived = new Derived();
        System.out.println(derived.value);
        System.out.println(derived.baseField);
        System.out.println(derived.rootField);
        System.out.println(derived.hidden);
        Base base = derived;
        System.out.println(base.hidden);
        System.out.println(base.hook());
        System.out.println(base.plain());
        return derived;
    }

    static void hiddenFields(Derived derived) {
        System.out.println(500);
        System.out.println(derived.hidden);
        Base base = derived;
        System.out.println(base.hidden);
        System.out.println(derived.derivedHidden());
        System.out.println(base.baseHidden());
        System.out.println(base.plain());
        System.out.println(derived.plain());
    }

    static void casts() {
        System.out.println(600);
        Object item = new Square(4);
        System.out.println(item instanceof Shape ? 1 : 0);
        Shape shape = (Shape) item;
        System.out.println(shape.area());
        Square square = (Square) shape;
        System.out.println(square.area());
        Object other = new Quiet();
        System.out.println(other instanceof Square ? 0 : 1);
        try {
            Square wrong = (Square) other;
            System.out.println(wrong.area());
        } catch (ClassCastException expected) {
            System.out.println(70);
        }
        Object nil = null;
        System.out.println(nil instanceof Square ? 0 : 1);
        Square none = (Square) nil;
        System.out.println(none == null ? 1 : 0);
        Animal cat = new Cat();
        System.out.println(cat instanceof Dog ? 0 : 1);
        System.out.println(item instanceof Animal ? 0 : 1);
    }

    static void factories() {
        System.out.println(700);
        Ticket ticket = Ticket.create(8);
        System.out.println(ticket.number());
        Ticket second = Ticket.create(9);
        System.out.println(second.number());
        PublicBox box = new PublicBox(10);
        System.out.println(box.payload());
        Holder holder = new Holder();
        holder.value = 7;
        System.out.println(holder.value);
    }

    public static void main(String[] args) {
        dispatch();
        defaults();
        abstractShapes();
        Derived derived = constructors();
        hiddenFields(derived);
        casts();
        factories();
    }
}
