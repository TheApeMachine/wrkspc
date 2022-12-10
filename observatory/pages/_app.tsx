import "../styles/globals.css"
import type { AppProps } from "next/app"
import Head from "next/head"
import { MantineProvider } from "@mantine/core"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
    <Head>
      <title>Page title</title>
      <meta name="viewport" content="minimum-scale=1, initial-scale=1, width=device-width" />
    </Head>
    <MantineProvider
      withGlobalStyles
      withNormalizeCSS
      theme={{
        globalStyles: (theme) => ({
          '*, *::before, *::after': {
            boxSizing: 'border-box',
          },

          html: {
            width: "100%", height: "100vh"
          },

          body: {
            ...theme.fn.fontStyles(),
            backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
            color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
            lineHeight: theme.lineHeight,
            width: "100%", height: "100vh"
          },
        }),
      }}
    >
      <Component style={{height: "100vh"}} {...pageProps} />
    </MantineProvider>
    </>
  )
}
