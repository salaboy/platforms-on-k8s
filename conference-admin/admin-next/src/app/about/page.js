import styles from '@/app/styles/about.module.css'

import Button from '../components/forms/button/button'
import Cloud from '../components/cloud/cloud'
import ExportedImage from "next-image-export-optimizer"
export default async function About() {
  
  
  return (
    <main className={styles.main}>
      
        <div className={`${styles.hero} ` }>
          <div className={ `grid content noMargin`}>
            <div className="col full">
              <h1>About <Cloud number="2" blue /></h1>
              
            </div>
          </div>
        </div>

      <div className="grid content">
      
        <div className="col third">
          
          
          <div>
            <h2>Repository</h2>
            <p>You can find the links to the source code and tutorials by going to the main Github repository: <a href="https://github.com/salaboy/platforms-on-k8s">github.com/salaboy/platforms-on-k8s</a></p>
            <br/>
            <Button main external link="https://github.com/salaboy/platforms-on-k8s">Go to the repository</Button>
          </div>
        </div>
          
            <div className="col third positionHalf">
              <div className={styles.book}>
                <div className={styles.bookContainer}>
                  <ExportedImage
                    src="/images/book.jpeg"
                    width={554}
                    height={694}
                    alt="Book"
                  />
                </div>
                <h2>Book</h2>
                <div >
                   <div>
                        <p>This application is fully covered by the <a target="_blank" href="http://mng.bz/jjKP">Plaform Engineering on Kubernetes Book</a>.</p>
                   </div>
                   
              
                </div>
            </div>
          </div>
          
      </div>
      
        <div className="grid content">
          <div className="col full">
            <h5>Developed and maintained by</h5>
            
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Mauricio Salatino</h3> 
                <p> <a href="https://www.salaboy.com" target={"_blank"}> salaboy.com </a></p>
                <ExportedImage
                  src="/images/salaboy.png"
                  width={100}
                  height={100}
                  alt="Salaboy"
                />
              </div>
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Ezequiel Salatino</h3>
                <p> <a href="https://salatino.me" target={"_blank"}> salatino.me </a></p>
                
                <ExportedImage
                  src="/images/esala.png"
                  width={100}
                  height={100}
                  alt="Esala"
                />
              </div>
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Matheus <br/> Cruz </h3>
                <p> <a href="https://twitter.com/mcruzdev1" target={"_blank"}> twitter.com/mcruzdev1 </a></p>
                
                <ExportedImage
                  src="/images/matheus.jpg"
                  width={100}
                  height={100}
                  alt="Matheus"
                />
              </div>
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Asare <br/> Nkansah </h3>
                <p> <a href="https://asarenkansah.github.io/asare-portfolio/" target={"_blank"}> asarenkansah.github.io/asare-portfolio/ </a></p>
                
                <ExportedImage
                  src="/images/asare.jpg"
                  width={100}
                  height={100}
                  alt="Asare"
                />
              </div>
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Marcos <br/> Lilljedahl </h3>
                <p> <a href="https://twitter.com/marcosnils" target={"_blank"}> twitter.com/marcosnils </a></p>
                
                <ExportedImage
                  src="/images/marcos.jpeg"
                  width={100}
                  height={100}
                  alt="Marcos"
                />
              </div>
            </div>
            <div className='col third'>
              <div className={styles.maintainer}>
                <h3>Giovanni <br/> Liva </h3>
                <p> <a href="https://twitter.com/thisthatDC" target={"_blank"}> twitter.com/thisthatDC </a></p>
                
                <ExportedImage
                  src="/images/giovanni.jpg"
                  width={100}
                  height={100}
                  alt="Giovani"
                />
              </div>
            </div>

            
          </div>
          <div className={styles.contribute}>
            <div className="grid content ">
              <div className="col half positionThird">
                <h5>Do you want to contribute to make this application better?</h5>
                <h3>
                 Go to the <a href="https://github.com/salaboy/platforms-on-k8s/issues">Platforms on K8s repository</a> and create an issue or drop me a message in Twitter <a href="https://twitter.com/salaboy">@Salaboy</a> 
                </h3>
                <br/>
                <Button  external link="https://twitter.com/salaboy">Drop me a message</Button>
              </div>
            </div>
          </div>
        
    </main>
  )
}
