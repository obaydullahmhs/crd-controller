# crd-controller
<html>
   <head>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
   </head>
   <body>
      <article id="0c4d8c73-c693-4414-aa53-27756e4904e6" class="page sans">
         <header>
            <h1 class="page-title">Sample Controller</h1>
         </header>
         <div class="page-body">
            <p id="2bb63bb3-4d85-4f14-a113-39e0817b3a32" class="">Clone this repository to this directory: <code>$GOPATH/src/obaydullahmhs</code></p>
            <p id="6e095295-5470-45f7-a60e-8219836c75c2" class="">To Clone Git with SSH: </p>
            <pre id="30609f9b-4e16-479b-b818-a1d9aca82055" class="code"><code>git clone git@github.com:obaydullahmhs/crd-controller.git</code></pre>
            <p id="2af473c3-b5d2-43aa-afa9-fbca65e53839" class="">Run this command to load/update package</p>
            <p id="8302f342-40cd-459d-8f6d-7b2425ef02a1" class=""> <code>go get all</code></p>
            <p id="edba26fd-c55d-41c4-882d-d979e7341870" class=""> <code>go mod tidy ; go mod vendor</code> </p>
            <p id="4a49a014-e2c7-4c72-8ad5-23696dc4a785" class=""></p>
            <p id="23276095-daf5-468a-92b9-c4d8bab6bac2" class="">Run the following command to give permission to some bash scripts</p>
            <p id="0ef15e16-9073-4b10-9f84-4e3f3e876e1f" class=""> <code>chmod +x hack/update-codegen.sh</code></p>
            <p id="83cf4ced-fbb4-4f55-a696-03496586c160" class=""> <code>chmod +x vendor/k8s.io/code-generator/generate-groups.sh</code> </p>
            <p id="897bd0d7-c519-4f58-b866-9ced84ed88dd" class=""></p>
            <p id="6412f1a3-a1b2-403e-9a88-a955af3ff265" class="">To generate clientset, listers and informers, run this command, this will also generate the custom resource definition yaml</p>
            <p id="f8bf07d4-98d6-46e5-9522-a570dfe598b7" class=""> <code>hack/update-codegen.sh</code> </p>
            <p id="7d80c32a-4943-4fab-a4c3-82dfe42bbec8" class=""></p>
            <p id="f251e647-d2af-4e67-92aa-5478b4bf238b" class="">To deploy your CRD in kubernetes cluster,</p>
            <p id="557e3401-1891-4cdf-a7c6-f6a156932cf1" class=""> <code>kubectl apply -f manifests/aadee.apps_aadees.yaml</code></p>
            <p id="8d7ca5ba-082b-440e-85d2-6ab3f1a441d1" class=""></p>
            <p id="e29eb70b-acc8-4a47-ab2e-f5e9873a551b" class="">Now to run the controller, build the project and run it using the following command:</p>
            <p id="a1f24ef8-658f-4274-8239-20db37c31d18" class=""><code>go build</code></p>
            <p id="aa322125-7b57-437e-ba1d-eac76dde03e1" class=""><code>./crd-controller</code></p>
            <p id="ccd6932d-ac80-43a7-86ad-dbdac8aa1dfb" class="">Now, all things are ready, You can deploy <code>aadee</code> type using <code>kubectl apply -f &lt;name&gt;.yaml</code> . Two demo yamls are in the manifests directory.</p>
            <p id="7519bc28-5888-4f31-932b-09e3f8c8901e" class="">This resource handles deployments type in kubernetes clusters.</p>
         </div>
      </article>
   </body>
</html>
