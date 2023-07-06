
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
                {pathname === "/backoffice/" && (
                    <>
                    <div className="col third">
                            <ul className={styles.logos}>
                                <li className={styles.logosItem} >
                                    <Link href="/"  className={pathname === "/" ? `${styles.active} ` : ' '} scroll={false}>
                                        <span className={styles.logosArrow}>
                                        <svg width="22" height="15" viewBox="0 0 22 15" fill="none" xmlns="http://www.w3.org/2000/svg">
<path fillRule="evenodd" clipRule="evenodd" d="M6.65685 0.292893L0.292892 6.65685C-0.0976315 7.04738 -0.0976315 7.68054 0.292892 8.07107L6.65685 14.435C7.04738 14.8256 7.68054 14.8256 8.07107 14.435C8.46159 14.0445 8.46159 13.4113 8.07107 13.0208L3.41421 8.36396H21C21.5523 8.36396 22 7.91625 22 7.36396C22 6.81168 21.5523 6.36396 21 6.36396H3.41421L8.07107 1.70711C8.46159 1.31658 8.46159 0.683418 8.07107 0.292893C7.68054 -0.0976311 7.04738 -0.0976311 6.65685 0.292893Z" fill="black"/>
</svg>

                                        </span>

                                        Back to website
                                    </Link>
                                </li>
                            </ul>
                        </div>
                    </>
                )}

                {pathname !== "/backoffice/" && (
                    <>
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
                                    <li className={`${styles.menuItem}   ${styles.icon} `}>
                                        <Link href="/backoffice/" className={pathname === "/backoffice" ? `${styles.active} ` : ' '} scroll={false}>
                                            <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <path d="M10 19C9.44772 19 9 19.4477 9 20C9 20.5523 9.44772 21 10 21V19ZM30.7071 20.7071C31.0976 20.3166 31.0976 19.6834 30.7071 19.2929L24.3431 12.9289C23.9526 12.5384 23.3195 12.5384 22.9289 12.9289C22.5384 13.3195 22.5384 13.9526 22.9289 14.3431L28.5858 20L22.9289 25.6569C22.5384 26.0474 22.5384 26.6805 22.9289 27.0711C23.3195 27.4616 23.9526 27.4616 24.3431 27.0711L30.7071 20.7071ZM10 21H30V19H10V21Z" fill="black"/>
                                            </svg>
                                        </Link></li>
                                </ul>
                            
                        </div>
                    </>
                )}
                
                
            </div>
        </nav>    
        
    
    );
}

