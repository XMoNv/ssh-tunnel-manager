import Api from './api.js'

export default {
    props: {
      data: Array,
      columns: Array,
      filterKey: String,
    },
    data() {
      return {
        sortKey: '',
        sortOrders: this.columns.reduce((o, key) => ((o[key] = 1), o), {})
      }
    },
    computed: {
      filteredData() {
        const sortKey = this.sortKey
        const filterKey = this.filterKey && this.filterKey.toLowerCase()
        const order = this.sortOrders[sortKey] || 1
        let data = this.data
        if (filterKey) {
          data = data.filter((row) => {
            return Object.keys(row).some((key) => {
              return String(row[key]).toLowerCase().indexOf(filterKey) > -1
            })
          })
        }
        if (sortKey) {
          data = data.slice().sort((a, b) => {
            a = a[sortKey]
            b = b[sortKey]
            return (a === b ? 0 : a > b ? 1 : -1) * order
          })
        }
        return data
      }
    },
    methods: {
      sortBy(key) {
        this.sortKey = key
        this.sortOrders[key] = this.sortOrders[key] * -1
      },
      capitalize(str) {
        return str.charAt(0).toUpperCase() + str.slice(1)
      },
      deleteTunnel(id) {
        Api.deleteTunnel(id).then((res) => {
            console.log(res)
            this.$parent.refresh()
        })
      },
      startTunnel(id) {
        Api.startTunnel(id).then((res) => {
            console.log(res)
            this.$parent.refresh()
        })
      },
      stopTunnel(id) {
        Api.stopTunnel(id).then((res) => {
            console.log(res)
            this.$parent.refresh()
        })
      },
    },
    template: `
    <table v-if="filteredData.length">
      <thead>
        <tr>
          <th v-for="key in columns"
            @click="sortBy(key)"
            :class="{ active: sortKey == key }">
            {{ capitalize(key) }}
            <span class="arrow" :class="sortOrders[key] > 0 ? 'asc' : 'dsc'">
            </span>
          </th>
          <th>Operation</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="entry in filteredData">
          <td v-for="key in columns">
            {{entry[key]}}
          </td>

          <td>
            <a v-if="entry.toggle" href="#" class="text-decoration-none" style="color:orange; padding:5px" @click="stopTunnel(entry.id)">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-stop-circle-fill" viewBox="0 0 16 16">
                <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0M6.5 5A1.5 1.5 0 0 0 5 6.5v3A1.5 1.5 0 0 0 6.5 11h3A1.5 1.5 0 0 0 11 9.5v-3A1.5 1.5 0 0 0 9.5 5z"></path>
                </svg>
            </a>
            <a v-else href="#" class="text-decoration-none" style="color:green; padding:5px"  @click="startTunnel(entry.id)">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-arrow-right-circle-fill" viewBox="0 0 16 16">
                <path d="M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0M4.5 7.5a.5.5 0 0 0 0 1h5.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 1 0-.708.708L10.293 7.5z"></path>
                </svg>
            </a>
            <a href="#" class="text-decoration-none" style="color:red; padding:5px" @click="deleteTunnel(entry.id)">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash-fill" viewBox="0 0 16 16">
                <path d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5M8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5m3 .5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 1 0"></path>
                </svg>
            </a>

          </td>
        </tr>
      </tbody>
    </table>
    <p v-else>No matches found.</p>
    `
  }
  