'use client'
import styles from './button.module.css'
import Link from 'next/link';

function Button({children, link, external, inline, clickHandler, small, main, disabled, inverted}) {

    var buttonElement;
    if(link){
      if(external){
        buttonElement = <a href={link} target="_blank">  <span>{children}</span> </a>
      }else {
        buttonElement = <Link href={link}>  <span>{children}</span> </Link>
      }
    }else {
      if(clickHandler){
        buttonElement = <Link className="__container" href="#" onClick={clickHandler}> <span>{children}</span>  </Link>
      }else {
        buttonElement = <Link href="#" className="__container">  <span>{children}</span> </Link>
      }
    }

    return (
      <div className={styles.button}>
            {buttonElement}
      </div>
    );

}
export default Button;