import type { NextPage } from 'next'
import Head from 'next/head'
import App from '../components/App'

const Home: NextPage = () => {
  return (
    <>
      <Head>
        <title>DTAP - Ragnarok Online Assistant</title>
      </Head>
      <App />
    </>
  )
}

export default Home
