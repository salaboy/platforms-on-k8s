'use client'
import styles from './button.module.css'
import Link from 'next/link';

function Button({children, link, external, inline, clickHandler, small, main, disabled, inverted, state}) {

    var buttonElement;
    if(link){
      if(external){
        buttonElement = <a href={link} target="_blank">  <span>{children}</span> </a>
      }else {
        buttonElement = <Link to={link}>  <span>{children}</span> </Link>
      }
    }else {
      if(clickHandler){
        buttonElement = <Link className="__container" href="#" onClick={clickHandler}> <span>{children}</span>  </Link>
      }else {
        buttonElement = <Link href="#" className="__container">  <span>{children}</span> </Link>
      }
    }

    return (
      <div>
            {buttonElement}
      </div>
    );

}
export default Button;