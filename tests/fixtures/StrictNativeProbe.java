/**
 * StrictNativeProbe — R2-A strict-native compliance test.
 *
 * Declares unresolved native methods of every return type and exercises
 * catch levels (UnsatisfiedLinkError, LinkageError, Error, Throwable),
 * finally execution, and independent-repeat behaviour.
 *
 * R2-A target (matches HotSpot):
 *   Every unresolved call throws a catchable UnsatisfiedLinkError. No zero
 *   value is ever produced. ULE, LinkageError, Error, and Throwable all catch.
 *   Finally always executes. Each call is an independent throw.
 */
public class StrictNativeProbe {
    static native void    nopVoid();
    static native boolean nopBool();
    static native byte    nopByte();
    static native char    nopChar();
    static native short   nopShort();
    static native int     nopInt();
    static native long    nopLong();
    static native float   nopFloat();
    static native double  nopDouble();
    static native String  nopString();
    static native int[]   nopIntArray();

    public static void main(String[] args) {
        int ules = 0;
        int total = 0;

        // --- int: catch via ULE exactly ---
        total++;
        try { nopInt(); System.out.println("FAIL int"); }
        catch (UnsatisfiedLinkError e) {
            ules++; System.out.println("ULE int ok");
        }

        // --- long: catch via ULE exactly ---
        total++;
        try { nopLong(); System.out.println("FAIL long"); }
        catch (UnsatisfiedLinkError e) {
            ules++; System.out.println("ULE long ok");
        }

        // --- String: catch via ULE exactly ---
        total++;
        try { nopString(); System.out.println("FAIL string"); }
        catch (UnsatisfiedLinkError e) {
            ules++; System.out.println("ULE string ok");
        }

        // --- void ---
        total++;
        try { nopVoid(); System.out.println("FAIL void"); }
        catch (UnsatisfiedLinkError e) {
            ules++; System.out.println("ULE void ok");
        }

        // --- boolean: finally runs even when ULE thrown ---
        total++;
        try { nopBool(); System.out.println("FAIL bool"); }
        catch (UnsatisfiedLinkError e) {
            ules++; System.out.println("ULE bool ok");
        } finally { System.out.println("finally ok"); }

        // --- Catch via LinkageError (supertype) — no ULE handler before it ---
        total++;
        try { nopInt(); System.out.println("FAIL LinkageError"); }
        catch (LinkageError e) { ules++; System.out.println("LinkageError ok"); }
        catch (Throwable t) { System.out.println("FAIL LinkageError wrong"); }

        // --- Catch via Error (supertype) — no LinkageError handler before it ---
        total++;
        try { nopInt(); System.out.println("FAIL Error"); }
        catch (Error e) { ules++; System.out.println("Error ok"); }
        catch (Throwable t) { System.out.println("FAIL Error wrong"); }

        // --- Catch via Throwable — no subtype handler before it ---
        total++;
        try { nopInt(); System.out.println("FAIL Throwable"); }
        catch (Throwable t) { ules++; System.out.println("Throwable ok"); }

        // --- Independent repeat calls ---
        total++;
        try { nopInt(); System.out.println("FAIL repeat1"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("repeat1 ok"); }

        total++;
        try { nopInt(); System.out.println("FAIL repeat2"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("repeat2 ok"); }

        // --- byte ---
        total++;
        try { nopByte(); System.out.println("FAIL byte"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE byte ok"); }

        // --- short ---
        total++;
        try { nopShort(); System.out.println("FAIL short"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE short ok"); }

        // --- char ---
        total++;
        try { nopChar(); System.out.println("FAIL char"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE char ok"); }

        // --- float ---
        total++;
        try { nopFloat(); System.out.println("FAIL float"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE float ok"); }

        // --- double ---
        total++;
        try { nopDouble(); System.out.println("FAIL double"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE double ok"); }

        // --- int[] ---
        total++;
        try { nopIntArray(); System.out.println("FAIL int[]"); }
        catch (UnsatisfiedLinkError e) { ules++; System.out.println("ULE int[] ok"); }

        // --- Summary ---
        if (ules == total) {
            System.out.print("PASS: ");
        } else {
            System.out.print("FAIL: ");
        }
        System.out.print(ules);
        System.out.print("/");
        System.out.print(total);
        System.out.println(" ULE caught");
    }
}
