<!-- 
The Hosts component presents an information table containing 
the mesh host info. It includes the hostname, number of CPUs, 
general top load report, RAM used and total RAM, the bandwidth 
available for experiments, the number of VMs, and host uptime.
 -->

<template>
  <div class="content">
    <hr>
    <b-table
      :data="hosts"
      :paginated="table.isPaginated"
      :per-page="table.perPage"
      :current-page.sync="table.currentPage"
      :pagination-simple="table.isPaginationSimple"
      :pagination-size="table.paginationSize"
      :default-sort-direction="table.defaultSortDirection"
      default-sort="name">
        <template slot-scope="props">
          <b-table-column field="name" label="Name" width="150" sortable>
            {{ hostName(props.row) }}
          </b-table-column>
          <b-table-column field="cpus" label="CPUs" width="50" sortable centered>
            {{ props.row.cpus }}
          </b-table-column>
          <b-table-column field="load" label="Load" width="250" centered>
            <span class="tag" :class="decorator(props.row.load[0], props.row.cpus)">
              {{ props.row.load[ 0 ] }}
            </span>
            --
            <span class="tag" :class="decorator(props.row.load[1], props.row.cpus)">
              {{ props.row.load[ 1 ] }}
            </span>
            --
            <span class="tag" :class="decorator(props.row.load[2], props.row.cpus)">
              {{ props.row.load[ 2 ] }}
            </span>
          </b-table-column>
          <b-table-column field="mem_used" label="RAM Used" width="100" centered>
            <span class="tag" :class="decorator(props.row.memused, props.row.memtotal)">
              {{ props.row.memused | ram }}
            </span>
          </b-table-column>
          <b-table-column field="mem_total" label="RAM Total" width="100" centered>
            {{ props.row.memtotal | ram }}
          </b-table-column>
          <b-table-column field="bandwidth" label="Bandwidth (MB/sec)" width="200" centered>
            {{ props.row.bandwidth }}
          </b-table-column>
          <b-table-column field="no_vms" label="# of VMs" width="100" sortable centered>
            {{ props.row.vms }}
          </b-table-column>
          <b-table-column field="uptime" label="Uptime" width="165">
            {{ props.row.uptime | uptime }}
          </b-table-column>
        </template>
    </b-table>
    <br>
    <b-field v-if="paginationNeeded" grouped position="is-right">
      <div class="control is-flex">
        <b-switch v-model="table.isPaginated" size="is-small" type="is-light" @input="changePaginate()">Paginate</b-switch>
      </div>
    </b-field>
    <b-loading :is-full-page="true" :active.sync="isWaiting" :can-cancel="false"></b-loading>
  </div>
</template>

<script>
  export default {
    
    beforeDestroy () {
      clearInterval( this.update );
    },
    
    created () {
      this.updateHosts();
      this.periodicUpdateHosts();
    },
    
    methods: {     
      updateHosts () {
        this.$http.get( 'hosts' ).then(
          response => {
            response.json().then(
              state => {
                if ( state.hosts.length == 0 ) {
                  this.isWaiting = true;
                } else {
                  this.hosts = state.hosts;
                  this.isWaiting = false;
                }
              }
            );
          }, response => {
            this.isWaiting = false;
            this.$buefy.toast.open({
              message: 'Getting the hosts failed with ' + response + ' response.',
              type: 'is-danger',
              duration: 4000
            });
          }
        );
      },

      periodicUpdateHosts () {
        this.update = setInterval( () => {
          this.updateHosts();
        }, 10000 )
      },

      changePaginate () {
        var user = localStorage.getItem( 'user' );
        localStorage.setItem( user + '.lastPaginate', this.table.isPaginated );
      },

      decorator ( sum, len ) {
        let actual = sum / len;
        if ( actual < .65 ) {
          return 'is-success';
        } else if ( actual >= .65 && actual < .85 ) {
          return 'is-warning';
        } else {
          return 'is-danger';
        }
      },

      hostName ( host ) {
        if ( host.headnode ) {
          return host.name + ' (headnode)';
        }

        return host.name;
      }
    },
    
    computed: {
      paginationNeeded () {
        var user = localStorage.getItem( 'user' );

        if ( localStorage.getItem( user + '.lastPaginate' ) ) {
          this.table.isPaginated = localStorage.getItem( user + '.lastPaginate' )  == 'true';
        }

        if ( this.hosts.length <= 10 ) {
          this.table.isPaginated = false;
          return false;
        } else {
          return true;
        }
      }
    },

    data () {
      return {
        table: {
          isPaginated: false,
          perPage: 10,
          currentPage: 1,
          isPaginationSimple: true,
          paginationSize: 'is-small',
          defaultSortDirection: 'asc'
        },
        hosts: [],
        isWaiting: true
      }
    }
  }
</script>
