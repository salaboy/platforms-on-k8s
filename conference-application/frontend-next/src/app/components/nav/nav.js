
'use client'
import styles from './nav.module.css'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import Cloud from '../cloud/cloud'



export default function Nav() {
    const pathname = usePathname()
    return (
        <nav className={styles.nav}>
            <div className="grid">
                <div className="col third">
                    <ul className={styles.logos}>
                        <li className={styles.logosItem} ><Link href="/"  className={pathname === "/" ? `${styles.active} ` : ' '} scroll={false}> <Cloud number="1" brand />CloudCon 2023</Link></li>
                    </ul>
                </div>
                <div className="col half positionHalf">
                    
                        <ul className={styles.menu}>
                            <li className={styles.menuItem}><Link href="/agenda/" className={pathname === "/agenda" ? `${styles.active} ` : ' '} scroll={false}>Agenda</Link></li>
                            <li className={styles.menuItem}><Link href="/proposals/" className={pathname === "/proposals" ? `${styles.active} ` : ' '} scroll={false}>Call for Proposals</Link></li>
                            <li className={styles.menuItem}><Link href="/about/" className={pathname === "/about" ? `${styles.active} ` : ' '} scroll={false}>About</Link></li>
                            <li className={styles.menuItem}><Link href="/backoffice/" className={pathname === "/backoffice" ? `${styles.active} ` : ' '} scroll={false}>Back Office</Link></li>
                        </ul>
                    
                </div>
                
            </div>
        </nav>    
        
    
    );
}

