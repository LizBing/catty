public class ProducerConsumer {
    static class OneSlot {
        int value;
        boolean full;

        synchronized void put(int next) throws InterruptedException {
            while (full) {
                wait();
            }
            value = next;
            full = true;
            notifyAll();
        }

        synchronized int take() throws InterruptedException {
            while (!full) {
                wait();
            }
            int result = value;
            full = false;
            notifyAll();
            return result;
        }
    }

    static final OneSlot SLOT = new OneSlot();
    static int sum;

    static class Producer extends Thread {
        public void run() {
            try {
                for (int i = 1; i <= 5; i++) {
                    SLOT.put(i);
                }
                SLOT.put(-1);
            } catch (InterruptedException unexpected) {
                throw new RuntimeException(unexpected);
            }
        }
    }

    static class Consumer extends Thread {
        public void run() {
            try {
                for (;;) {
                    int value = SLOT.take();
                    if (value == -1) {
                        return;
                    }
                    sum += value;
                }
            } catch (InterruptedException unexpected) {
                throw new RuntimeException(unexpected);
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Producer producer = new Producer();
        Consumer consumer = new Consumer();
        consumer.start();
        producer.start();
        producer.join();
        consumer.join();
        System.out.println(sum);
    }
}
