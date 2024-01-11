<script>
  import { onMount } from 'svelte';
  import { getAddress } from '@ethersproject/address';
  import { formatEther } from '@ethersproject/units';
  import { setDefaults as setToast, toast } from 'bulma-toast';
    import App from './App.svelte';

  $: address = null
  $: network = null
  $: disabled = false
  $: faucetInfo = {
    account: '0x0000000000000000000000000000000000000000',
    network: 'Testnet',
    payout: 1,
    rpcAvailable: true,
  }
  const testnetConfig = {
    chainId: '0xb1ce2',
    chainName: 'OnlyLayer',
    nativeCurrency: {
      name: 'OnlyLayer',
      symbol: 'ONLY',
      decimals: 18,
    },
    rpcUrls: [
      'https://onlylayer.org/',
    ],
    blockExplorerUrls: [
      'https://onlyscan.info/',
    ],
    iconUrls: [],
  }
  const updateChainId = async () => {
    if (!ethereum) {
      return
    }
    const chainId = await ethereum.request({
      method: 'eth_chainId',
    })
    network = chainId
  }
  let address = null
  let network = null

  const requestAccounts = async () => {
    if (!ethereum) {
      return
    }
    const accounts = await ethereum.request({
      method: 'eth_requestAccounts',
    })
    address = accounts[0]
  }

  const buttonfunction = async() => {
    if (addTestnetToMetamask) {
      if (!ethereum) {
        return
      }
      network = ethereum.chainId
      const requestAccountsOnce = (canRun = true) => () => {
        // ethereum.off('connect', requestAccountsOnce)
        if (!canRun) {
          return
        }
        requestAccounts()
        canRun = false
        document.querySelector('.button').innerText = 'Wallet Connected'
      }
      ethereum.on('connect', requestAccountsOnce)
      ethereum.on('connect', updateChainId)
      ethereum.on('chainChanged', updateChainId)
      ethereum.on('accountsChanged', (accounts) => {
        address = accounts[0]
      })
    }
  }
  
  buttonfunction()
  onMount(async () => {
    const res = await fetch('/api/info');
    try {
      faucetInfo = await res.json();
    } catch (err) {}
    if (!faucetInfo.account || faucetInfo.account === '0x0000000000000000000000000000000000000000') {
      network = testnetConfig.chainId
      disabled = true
      return
    }
    updateChainId()
    faucetInfo.network = capitalize(faucetInfo.network);
    faucetInfo.payout = parseInt(formatEther(faucetInfo.payout));
  });

  setToast({
    duration: 10000,
    position: 'bottom-center',
    dismissible: true,
    pauseOnHover: true,
    closeOnClick: false,
    animate: { in: 'fadeIn', out: 'fadeOut' },
  });

  async function addTestnetToMetamask () {
    try {
      await ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: testnetConfig.chainId }],
      })
    } catch (switchError) {
      if (switchError.code === 4902) {
        try {
          await ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [testnetConfig],
          })
        } catch (addError) {
          console.log('unable to add a new chain', addError)
          return
        }
      }
      console.log('unable to complete request', switchError)
    }
  }

  async function handleRequest() {
    try {
      address = getAddress(address);
    } catch (error) {
      toast({ message: error.reason, type: 'is-warning' });
      return;
    }

    let formData = new FormData();
    formData.append('address', address);
    const res = await fetch('/api/claim', {
      method: 'POST',
      body: formData,
    })
    const { ok } = res
    let type = ok ? 'is-success' : 'is-warning'
    let message = await res.text()
    toast({ message, type });
  }

  function capitalize(str) {
    const lower = str.toLowerCase();
    return str.charAt(0).toUpperCase() + lower.slice(1);
  }
</script>

<svelte:head>
  <title>OnlyLayer Testnet Faucet</title>
</svelte:head>

<main>
  <section class="hero is-info is-fullheight">
    <div class="hero-head">
      <nav class="navbar">
        <div class="container">
          <div class="navbar-brand">
            <a class="navbar-item" href=".">
              <img class="navbar-logo" src="./favicon-32x32.png" alt="" />
              <span><b>OnlyLayer Faucet</b></span>
              <span class="icon">
                <i class="fa fa-bath" />
              </span>
            </a>
          </div>
          <div id="navbarMenu" class="navbar-menu">
            <div class="navbar-end">
              <span class="navbar-item">
              <button on:click={requestAccounts} class="button is-primary is-rounded short-text">{address ? 'Wallet Connected' : 'Connect Wallet'}</button>
            </div>
          </div>
        </div>
      </nav>
    </div>

    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-6 is-offset-3">
          <h1 class="title">
            Receive {faucetInfo.payout} coins per request
          </h1>
          <h2 class="subtitle">
            Serving from
            {faucetInfo.account}
          </h2>
          {#if disabled}
            <div class="box">
              <p>unable to contact network. please try again later</p>
            </div>
          {:else if network === testnetConfig.chainId}
          <div class="box">
            <div class="field is-grouped">
              <p class="control is-expanded">
                <input
                  bind:value={address}
                  class="input is-rounded"
                  type="text"
                  placeholder="Enter your address"
                />
              </p>
              <p class="control">
                <button
                  on:click={handleRequest}
                  class="button is-primary is-rounded"
                >
                  Request
                </button>
              </p>
            </div>
          </div>
          {:else}
          <button
            on:click={addTestnetToMetamask}
            class="button is-primary is-rounded">Switch to / Add OnlyLayer Testnet.</button>
          {/if}
        </div>
      </div>
    </div>
  </section>
</main>

<style>
  .hero.is-info {
    background: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)),
      url('/bg.jpg') no-repeat center center fixed;
    -webkit-background-size: cover;
    -moz-background-size: cover;
    -o-background-size: cover;
    background-size: cover;
  }
  @media only screen and (max-width: 767px) {
    .short-text {
      font-size: 14px;
      line-height: 1.2;
    }
  }  
  .hero .subtitle {
    padding: 3rem 0;
    line-height: 1.5;
  }
  .box {
    border-radius: 19px;
  }

  .navbar-logo {
    margin-right: 10px;
    height: 32px; 
    width: 32px; 
  }
</style>
