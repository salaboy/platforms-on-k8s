
'use client'
import styles from './footer.module.css'
import Link from 'next/link'
import { usePathname } from 'next/navigation'



export default function Footer() {
    const pathname = usePathname()
    return (
        <nav className={styles.footer}>
            <div className="grid content noMargin">
                <div className="col third">
                    <ul className={styles.logos}>
                        <li className={styles.logosItem} ><h4><Link href="/"  className={pathname === "/" ? `${styles.active} ` : ' '} scroll={false}>Cloud-Native <br /> Conf 2023.</Link></h4></li>
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
                            
                            <li className={styles.menuItem}><Link href="/backoffice/" className={pathname === "/backoffice" ? `${styles.active} ` : ' '} scroll={false}>Go to Backoffice</Link></li>
                        </ul>
                        <p className='p --s'>
                            Copyright Salaboy. 2023. <br /> Visit <a href="https://salaboy.com" target='_blank'>my blog</a> for more information.
                        </p>
                        
                </div>
                
            </div>
        </nav>    
        
    
    );
}

