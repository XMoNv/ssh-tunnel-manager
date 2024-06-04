
const API_URL =  "/api/v1/"
const Api = {

    methods: {
        async post(url = "", data = {}) {
            // Default options are marked with *
            const response = await fetch(url, {
              method: "POST", // *GET, POST, PUT, DELETE, etc.
              mode: "cors", // no-cors, *cors, same-origin
              cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
              credentials: "same-origin", // include, *same-origin, omit
              headers: {
                "Content-Type": "application/json",
                // 'Content-Type': 'application/x-www-form-urlencoded',
              },
              redirect: "follow", // manual, *follow, error
              referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
              body: JSON.stringify(data), // body data type must match "Content-Type" header
            });
            return response.json(); // parses JSON response into native JavaScript objects
        },

        async fetchConfigList() {
            const url = `${API_URL}list`
            console.log('request', url)
            var result = await (await fetch(url)).json()
            return result['configs']
        },

        async createTunnel(d) {
            const url = `${API_URL}create`
            console.log('request', url)
            var result = await (await this.post(url, {
                "name": d.name,
                "mode": ">",
                "user": d.user,
                "host": d.host,
                "port": d.port,
                "bindAddr": d.bindAddr,
                "dialAddr": d.dialAddr,
                "localPort": d.localPort,
                "remotePort": d.remotePort,
                "toggle": d.toggle ? 1 : 0
            }))
            return result.code
        },

        async descTunnel(id) {
            const url = `${API_URL}info/${id}`
            console.log('request', url)
            var result = await (await fetch(url))
            return result
        },

        async startTunnel(id) {
            const url = `${API_URL}start/${id}`
            console.log('request', url)
            var result = await (await this.post(url, {}))
            return result.code
        },

        async stopTunnel(id) {
            const url = `${API_URL}stop/${id}`
            console.log('request', url)
            var result = await (await this.post(url, {}))
            return result.code
        },

        async deleteTunnel(id) {
            const url = `${API_URL}delete/${id}`
            console.log('request', url)
            var result = await (await this.post(url, {}))
            return result.code
        },
    }
  }
  

  export default Api.methods