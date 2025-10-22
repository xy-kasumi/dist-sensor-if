const { createApp } = Vue;

createApp({
    data() {
        return {
            distance: null,
            error: null,
            loading: false
        }
    },
    methods: {
        async readDistance() {
            this.loading = true;
            this.error = null;
            this.distance = null;

            try {
                const response = await fetch('/api/command', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ command: 0xb001 })
                });

                const result = await response.json();

                if (result.ok) {
                    const sint16 = result.data > 32767 ? result.data - 65536 : result.data;
                    this.distance = sint16;
                } else {
                    this.error = result.error || 'Unknown error';
                }
            } catch (err) {
                this.error = `Network error: ${err.message}`;
            } finally {
                this.loading = false;
            }
        }
    }
}).mount('#app');
