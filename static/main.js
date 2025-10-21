const { createApp } = Vue;

createApp({
    data() {
        return {
            message: 'Hello Vue 3!'
        }
    },
    template: `
        <div>
            <h1>{{ message }}</h1>
        </div>
    `
}).mount('#app');
