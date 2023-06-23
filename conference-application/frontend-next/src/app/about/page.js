import styles from '@/app/styles/about.module.css'

import Button from '../components/forms/button/button'
import Image from 'next/image'

export default async function About() {
  
  
  return (
    <main className={styles.main}>
      
        <div className={`${styles.hero} ` }>
          <div className={ `grid content noMargin`}>
            <div className="col full">
              <h1>About</h1>
              
            </div>
          </div>
        </div>

      <div className="grid content">
      
        <div className="col full">
          
          <div className="grid">
          <div>
            <h4>Repository</h4>
            <p>You can find the links to the source code and tutorials by going to the main Github repository: <a href="https://github.com/salaboy/from-monolith-to-k8s">https://github.com/salaboy/from-monolith-to-k8s</a></p>
            <br/>
            <Button main external link="https://github.com/salaboy/platforms-on-k8s">Go to the repository</Button>
          </div>
          <div>
            <h4>Book</h4>
            <div className="grid">
            <div>
              <p>This application is fully covered by the <a target="_blank" href="http://mng.bz/jjKP">Plaform Engineering on Kubernetes Book</a>.</p>
            </div>
            <div>
            <Image
              src="/images/book.jpeg"
              width={554}
              height={694}
              alt="Book"
            />
            </div>
            </div>
          </div>
          </div>
        </div>
      
        <div className="grid">
          <div>
            <h4>Developed and maintained by</h4>
            https://twitter.com/salaboy
            <Image
              src="/images/salaboy.png"
              width={100}
              height={100}
              alt="Salaboy"
            />
            https://salatino.me
            <Image
              src="/images/esala.png"
              width={100}
              height={100}
              alt="Esala"
            />
            https://twitter.com/mcruzdev1
            <Image
              src="/images/matheus.jpg"
              width={100}
              height={100}
              alt="Matheus"
            />
       
          
          </div>
          <div>
          <h4>Contribute</h4>
            <p className="p p-b">Do you want to contribute to make this application better?
            Go to the <a href="https://github.com/salaboy/platforms-on-k8s/issues">Platforms on K8s repository</a> and create an issue or drop me a message in Twitter <a href="https://twitter.com/salaboy">@Salaboy</a> </p>
            <br/>
            <Button  external link="https://twitter.com/salaboy">Drop me a message</Button>
          </div>
        </div> 
        
        
      </div> 
      

       
    </main>
  )
}
