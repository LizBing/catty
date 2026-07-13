public class FinalFieldPublication {
    static class Holder {
        final int finalValue;
        int ordinaryValue;

        Holder() {
            finalValue = 42;
            ordinaryValue = 24;
        }
    }

    static volatile Holder shared;
    static int observedFinal;
    static int observedOrdinary;

    static class Worker extends Thread {
        public void run() {
            Holder holder;
            while ((holder = shared) == null) {
                Thread.onSpinWait();
            }
            observedFinal = holder.finalValue;
            observedOrdinary = holder.ordinaryValue;
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        shared = new Holder();
        worker.join();
        System.out.println(observedFinal);
        System.out.println(observedOrdinary);
    }
}
