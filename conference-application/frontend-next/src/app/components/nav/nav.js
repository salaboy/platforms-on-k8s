
import styles from './nav.module.css'
import Link from 'next/link'
import { useRouter } from 'next/navigation';



export default function Nav() {
    
    
    return (
        <>
        <nav className={styles.nav}>
            <ul className={styles.menu}>
                <li className={styles.menuItem}><Link href="/" scroll={false}>Home</Link></li>
                <li className={styles.menuItem}><Link href="/about/" scroll={false}>About</Link></li>
                <li className={styles.menuItem}><Link href="/agenda/" scroll={false}>Agenda</Link></li>
                <li className={styles.menuItem}><Link href="/proposals/" scroll={false}>Proposals</Link></li>
                <li className={styles.menuItem}><Link href="/backend/" scroll={false}>Backend</Link></li>
            </ul>
        </nav>
            
     
                    
            
        </>
        
    
    );
}