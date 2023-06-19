
'use client'
import styles from './footer.module.css'
import Link from 'next/link'
import { usePathname } from 'next/navigation'



export default function Footer() {
    const pathname = usePathname()
    return (
        <nav className={styles.footer}>
            <div className="grid">
                <div className="col third">
                    <ul className={styles.logos}>
                        <li className={styles.logosItem} ><Link href="/"  className={pathname === "/" ? `${styles.active} ` : ' '} scroll={false}>CloudCon 2023.</Link></li>
                    </ul>
                </div>
                <div className="col third ">
                    
                        <ul className={styles.menu}>
                            <li className={styles.menuItem}><Link href="/about/" className={pathname === "/about" ? `${styles.active} ` : ' '} scroll={false}>About</Link></li>
                            <li className={styles.menuItem}><Link href="/agenda/" className={pathname === "/agenda" ? `${styles.active} ` : ' '} scroll={false}>Agenda</Link></li>
                            <li className={styles.menuItem}><Link href="/proposals/" className={pathname === "/proposals" ? `${styles.active} ` : ' '} scroll={false}>Proposals</Link></li>
                            
                        </ul>
                    
                </div>
                <div className="col third ">
                    
                        <ul className={styles.menu}>
                            
                            <li className={styles.menuItem}><Link href="/backend/" className={pathname === "/backend" ? `${styles.active} ` : ' '} scroll={false}>Backend</Link></li>
                        </ul>
                        
                </div>
                
            </div>
        </nav>    
        
    
    );
}

