import styles from "../styles/Home.module.css"
import Link from "next/link"

export default function Home() {
  return (
    <div className={styles.container}>
      <Link href="/builder/editor">Builder</Link>
    </div>
  )
}
