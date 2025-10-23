const { createApp } = Vue;

/**
 * Send command.
 * @param {number} cmd 
 * @returns {Promise<number>} data
 */
const command = async (cmd) => {
    const response = await fetch('/api/command', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ command: cmd })
    });
    const result = await response.json();
    if (!result.ok) {
        throw new Error(result.error || 'Unknown error');
    }
    return result.data;
};

/**
 * Measure distance. Throws if error.
 * @returns {Promise<number | null>} distance in mm or null if out of range.
 */
const measureDist = async () => {
    const data = await command(0xb001);
    const distUm = data > 32767 ? data - 65536 : data;
    return distUm > 10000 ? null : distUm * 1e-3;
};

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

createApp({
    data() {
        return {
            error: null,
            distance: null,
            minDist: null,
            maxDist: null,
            running: false,
        }
    },
    async mounted() {
        while (true) {
            try {
                const dist = await measureDist();
                if (dist === null) {
                    this.error = 'Out of range';
                    this.distance = null;
                } else {
                    this.error = null;
                    this.distance = dist;
                    if (this.running) {
                        if (this.minDist === null || dist < this.minDist) {
                            this.minDist = dist;
                        }
                        if (this.maxDist === null || dist > this.maxDist) {
                            this.maxDist = dist;
                        }
                    }
                }
            } catch (error) {
                this.error = `Device error: ${error.message}`;
            }
            await sleep(50);
        }
    },
    computed: {
        delta() {
            if (this.minDist === null || this.maxDist === null) {
                return null;
            }
            return this.maxDist - this.minDist;
        },
    },
    methods: {
        start() {
            this.minDist = null;
            this.maxDist = null;
            this.running = true;
        },
        stop() {
            this.running = false;
        },
    }
}).mount('#app');
