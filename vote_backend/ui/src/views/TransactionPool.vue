<template>
  <div class="">
    <div class="ml-96 mr-12 mt-32 flex flex-col justify-start  bg-[#ffffff] h-full">
      <div class="flex flex-row justify-between items-center w-full">
        <div>
          <p class="px-12">Transactions</p>
        </div>

        <div class="flex flex-row px-12">
          <router-link :to="{ name: 'transaction-pool' }" @click="$router.go()" class="flex items-center p-2">
            <img src="../assets/images/refresh.svg" class="w-6 h-8">
          </router-link>

        </div>
      </div>
      <div class="px-12 flex justify-center w-fit">

        <table class="table-auto space-y-6 w-fit">
          <thead class="">

            <tr class="rounded-3xl">
              <th class="px-1.5 py-2 border">Index</th>
              <th class="px-1.5 py-2 border">Transaction ID</th>
              <th class="px-1.5 py-2 border">Node ID</th>
              <th class="px-1.5 py-2 border">Candidate</th>
              <th class="px-1.5 py-2 border">County</th>
              <th class="px-1.5 py-2 border">Constituency</th>
              <th class="px-1.5 py-2 border">Ward</th>
              <th class="px-1.5 py-2 border">Polling Station</th>

            </tr>
          </thead>
          <tbody>
            <tr v-for="(transaction, index) in this.allTransactions" :key="index" class="bg-white border-b">
              <td class="px-1.5 py-2   border ">{{ index + 1 }} </td>
              <td class="px-1.5 py-2   border break-all">{{ transaction.Txid }}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.NodeId }}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.Candidate}}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.County }}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.Constituency }}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.Ward }}</td>
              <td class="px-1.5 py-2   border ">{{ transaction.PollingStation }}</td>

            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import SecureLS from 'secure-ls'

export default {

  data () {
    this.getTransactions()
    return {
      allTransactions: {}

    }
  },

  methods: {
    getTransactions () {
      const ls = new SecureLS()
      const config = {
        headers: { Authorization: `Bearer ${ls.get('user').token}` }
      }
      axios.get(
        'http://127.0.0.1:3500/api/secured/get-all-transactions', config
      ).then((response) => {
        console.log(response.data)
        this.allTransactions = response.data
      }).catch(function (error) {
        if (error.response.status === 401) {
          ls.removeAll()
          window.location.href = '/'
        } if (error.toJSON().message === 'Network Error') {
          alert('no internet connection')
        }
      })
    }
  }

}
</script>

<style scoped></style>
