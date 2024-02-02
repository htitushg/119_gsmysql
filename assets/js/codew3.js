function demarrer(){
    document.getElementById("bouton").addEventListener("click", affichemdp);
}
function affichemdp(){
    if (document.getElementById("passid").type=="text")
    {
        document.getElementById("passid").type="password";
        document.getElementById("passid2").type="password";
    }
    else
    {
        document.getElementById("passid").type="text";
        document.getElementById("passid2").type="text";
    }
}
//Corps
//Attends le chargement de la page pour exécuter la fonction demarrer
window.addEventListener("load", demarrer);
//____________________________________________
function onlycheck(checkbox) {
    var checkboxes = document.getElementsByName('check');
    checkboxes.forEach((item) => {
        if (item !== checkbox) item.checked = false;
    });
}
//____________________________________________
function annulation()
{
  alert('Annulation demandée');
      window.location.href("/");
      return true;
}
//____________________________________________
function formValidation()
{
  //var resultat = true
  var uid = document.registration.pseudo;
  var uemail = document.registration.email;
  var passid = document.registration.passid;
  var passid2 = document.registration.passid2;
  var fname = document.registration.firstname;
  var lname = document.registration.lastname;
  var uadd = document.registration.address;
  var town = document.registration.town;
  var uzip = document.registration.zip;
  var ucountry = document.registration.country;
  var umsex = document.registration.Homme;
  var ufsex = document.registration.Femme; 
  if(pseudo_validation(uid,5,12))
  {if(validateEmail(uemail))
   {if(passid_validation(passid,passid2,8,20))
    {if(allLetter(fname,"fname"))
     {if(allLetter(lname,"lname"))
      {if(alphanumeric(uadd))
       {if(allnumeric(uzip))
        {if(allLetter(town,"town"))
         {if(countryselect(ucountry))
          {if(validateSex(umsex,ufsex))
           {           
           }
          }
         }
        }    
       }
      }
     }
    }
   }   
  }
  return false;
}

//______________________________________________________________________________________________

function pseudo_validation(uid,mx,my)
{
  var letters = /^[A-Za-z0-9]+$/;
  if(uid.value.match(letters))
  {
    //alert("Pseudo correct")
    return true;
  }
  else
  {
    alert("Pseudo must have alphabet characters and digits");
    uid.focus();
    return false;
  }
}
//______________________________________________________________________________________________

function passid_validation(passid,passid2,mx,my)
{
 var decimal=  /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z0-9])(?!.*\s).{8,20}$/;
 if(passid.value.match(decimal)) 
 { 
    //le mot de passe est conforme
    if(passid.value == passid2.value)
    {
        //alert("Password Correct, les 2 entrées sont bien conformes et identiques");
        return true;
    }
    else
    {
        alert("Password Incorrect, les 2 entrées doivent être conformes et identiques"); 
        passid.focus();
        return false;
    }
 }
 else
 { 
    alert("Le mot de passe doit avoir 8 à 20 caractères, au moins 1 chiffre, 1 majuscule, 1 minuscule, 1 caractère spécial !");
    passid.focus();
    return false;
 }
}
//______________________________________________________________________________________________

function allLetter(uname, nomchamp)
{ 
  //var letters = /^[a-zA-Z -]+$/;
  var letters = /^[A-Za-z -']+$/;
  if(uname.value.match(letters))
  {
    //alert("Pseudo correct")
    return true;
  }
  else
  {
    alert(nomchamp +" must have alphabet characters only");
    uname.focus();
    return false;
  }
}
//______________________________________________________________________________________________

function alphanumeric(uadd)
{ 
  var letters = /^[0-9a-zA-Z'éèàùë\- ]+$/;
  //var letters = /^[0-9a-zA-Z]+$/;
  if(uadd.value.match(letters)) 
  {
    //alert("Adress correct");
    return true;
  }
  else
  {
    alert("User address must have alphanumeric characters only");
    uadd.focus();
    return false;
  }
}
//______________________________________________________________________________________________

function countryselect(ucountry)
{
  if(ucountry.value == "Default")
  {
    alert("Select your country from the list");
    ucountry.focus();
    //resultat= false;
    return false;
  }
  else
  {
    //alert("Country correct");
    return true;
  }
}
//______________________________________________________________________________________________

function allnumeric(uzip)
{ 
  var numbers = /^[0-9]+$/;
  if(uzip.value.match(numbers))
  {
    //alert("Zip correct");
    return true;
  }
  else
  {
    alert("ZIP code must have numeric characters only");
    uzip.focus();
    return false;
  }
}
//______________________________________________________________________________________________

function validateEmail(uemail)
{
  var mailformat = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
  if(uemail.value.match(mailformat))
  {
    //alert("Email correct");
    return true;
  }
  else
  {
    alert("You have entered an invalid email address!");
    uemail.focus();
    return false;
  }
}
function validateSex(umsex,ufsex)
{
    //alert("Vérification Selection Masculin/Féminin");
    if(((umsex.checked)&&((ufsex.checked)))||((!umsex.checked)&&((!ufsex.checked)))) 
    {
        alert("Select Masculin/Féminin");
        umsex.focus();
        return false;
    } 
    else
    {
      alert('Form Successfully Submitted');
      window.location.href("/RegisterPost");
      return true;
    }
}


