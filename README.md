#  Cloud Provisioner CLI

Un outil en ligne de commande (CLI) développé en Go pour automatiser le provisionnement, la gestion du statut et la destruction de machines virtuelles sur **Amazon Web Services (AWS)** via le **SDK AWS v2**.

Ce projet suit les principes de la **Clean Architecture** (conception *cloud-agnostic* basée sur des interfaces) et intègre la gestion synchrone des états via les **Waiters AWS** et les **Contexte Go (`context.Context`)**.

---

##  Architecture du Projet

Le projet est structuré selon la convention **Standard Go Project Layout** :

```text
cloud-provisioner/
├── cmd/
│   └── main.go           # Point d'entrée de la CLI (routing & flags)
├── pkg/
│   ├── cloud/
│   │   ├── provider.go   # Interface CloudProvider et structures génériques
│   │   └── aws.go        # Implémentation du SDK AWS v2 (EC2, Waiters)
│   └── config/
│       └── config.go     # Chargement de la configuration & variables d'environnement
├── go.mod                # Dépendances du module Go
└── README.md             # Documentation du projet

```

---

##  Fonctionnalités

* **Architecture découplée** : Interface `CloudProvider` permettant de remplacer ou d'ajouter facilement d'autres fournisseurs Cloud (GCP, Azure, LocalStack).
* **Sous-commandes CLI** : `create`, `status`, et `destroy`.
* **Démarrage synchrone** : Attente active (`InstanceRunningWaiter`) jusqu'à ce que la VM soit complètement opérationnelle et qu'une IP publique lui soit attribuée.
* **Gestion des Timeouts** : Utilisation systématique de `context.WithTimeout` pour éviter tout blocage indéfini lors des appels réseau/API.
* **Configuration flexible** : Prise en charge des variables d'environnement avec des valeurs de secours (*fallbacks*).

---

##  Prérequis

* **Go** : Version `1.20` ou supérieure.
* **Compte AWS** : Avec des clés d'accès IAM valides (possédant les droits `ec2:RunInstances`, `ec2:DescribeInstances`, `ec2:TerminateInstances`).

---

## ⚙️ Configuration & Identifiants AWS

L'application lit automatiquement vos identifiants AWS depuis votre environnement système ou votre fichier local `~/.aws/credentials`.

Exportez vos identifiants dans votre terminal :

```bash
export AWS_ACCESS_KEY_ID="VOTRE_ACCESS_KEY"
export AWS_SECRET_ACCESS_KEY="VOTRE_SECRET_KEY"

```

Vous pouvez également personnaliser le comportement par défaut via les variables d'environnement suivantes :

| Variable | Description | Valeur par défaut |
| --- | --- | --- |
| `AWS_REGION` | Région AWS cible | `eu-west-3` (Paris) |
| `DEFAULT_INSTANCE_TYPE` | Type d'instance EC2 par défaut | `t2.micro` |
| `DEFAULT_AMI_ID` | Identifiant d'AMI par défaut | `""` (Requis pour la création) |

---

##  Compilation & Installation

À la racine du projet, compilez le binaire exécutable :

```bash
# Compilation du binaire
go build -o infra cmd/main.go

```

Un exécutable nommé `infra` (ou `infra.exe` sous Windows) sera généré.

---

##  Guide d'Utilisation

### 1. Provisionner une nouvelle instance VM (`create`)

Crée une instance EC2 et bloque le terminal jusqu'à ce qu'elle soit en état `running` et dotée d'une IP publique.

```bash
./infra create -name=prod-web-server -ami=ami-00c7042a420b784a6 -type=t2.micro

```

**Options :**

* `-name` : Nom attribué à la ressource via un Tag EC2 (Défaut : `go-vm`).
* `-type` : Type d'instance EC2 (Défaut : `t2.micro`).
* `-ami` : ID de l'image système (Obligatoire sauf si défini via `DEFAULT_AMI_ID`).

---

### 2. Consulter le statut d'une instance (`status`)

Récupère l'état actuel (`pending`, `running`, `terminated`), l'IP publique et l'IP privée d'une VM via son ID.

```bash
./infra status -id=i-0123456789abcdef0

```

---

### 3. Résilier/Détruire une instance (`destroy`)

Envoie un ordre de suppression définitif à AWS pour stopper et détruire la VM.

```bash
./infra destroy -id=i-0123456789abcdef0

```

---

##  Exemple de Sortie Terminal (`create`)

```text
 Initialisation du Cloud Provisioner...
 Création de l'instance 'prod-web-server' (t2.micro) dans la région eu-west-3...
 Instance créée avec succès ! ID: i-0abcd1234efgh5678 (État: pending)
 Attente du démarrage complet de l'instance i-0abcd1234efgh5678...

 --- INSTANCE PRÊTE A L'EMPLOI ---
 ID         : i-0abcd1234efgh5678
 État       : running
 IP Publique : 13.38.120.45
 IP Privée   : 172.31.22.105

```
