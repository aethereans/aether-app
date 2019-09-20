<template>
  <div class="onboard-sublocation">
    <div class="onboard-carrier">
      <transition name="fade" appear>
        <div class="conditions-container">
          <a-markdown class="onboard-markdown" :content="content"></a-markdown>
          <div class="tocbox">
            <a-markdown class="onboard-markdown" :content="toc"></a-markdown>
          </div>
        </div>
      </transition>
      <div class="continue-box">
        <router-link class="button is-outlined next-button" to="/onboard/6">
          AGREE
        </router-link>
        <div class="button is-outlined next-button" @click="quitApp">
          DECLINE AND QUIT
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
var ipc = require('../../../../../node_modules/electron-better-ipc')
export default {
  name: 'onboard5',
  data() {
    return {
      content: textContent,
      toc: termsAndConditions,
    }
  },
  methods: {
    quitApp(this: any) {
      ipc.callMain('QuitApp')
    },
  },
}
var textContent = `
## *Information wants to be free*

If you post anything on Aether, you agree that your posts will be licensed under Creative Commons [BY-SA](https://creativecommons.org/licenses/by-sa/4.0/) license.

This is so that everyone can share content you create.

You also agree to the [Terms and Conditions](https://getaether.net/terms-and-conditions/), provided below.
  `
var termsAndConditions = `
Agreement between User, and Getaether.net, Aether App

The Getaether.net website (the "Site") and Aether app (the "App") is comprised of various web and application pages operated by Aether Technologies, Inc. ("Aether Technologies Inc."). Getaether.net and Aether app is offered to you conditioned on your acceptance without modification of the terms, conditions, and notices contained herein (the "Terms"). Your use of Getaether.net website and Aether app constitutes your agreement to all such Terms. Please read these terms carefully, and keep a copy of them for your reference.


Getaether.net is a Forum/Online discussion Site. Aether App is a Forum/Online discussion App.


Specifically, Aether is an app that allows its users to read, and participate in peer-to-peer, decentralized online discussion forums.


Privacy

Your use of Getaether.net and Aether App is subject to Aether Technologies Inc.'s Privacy Policy. Please review our Privacy Policy, which also governs the Site and App, and informs users of our data collection practices.

Electronic Communications

Visiting Getaether.net, using the Aether App, or sending emails to Aether Technologies Inc. constitutes electronic communications. You consent to receive electronic communications and you agree that all agreements, notices, disclosures and other communications that we provide to you electronically, via email and on the Site or App, satisfy any legal requirement that such communications be in writing.

Your Account

If you use this site or App, you are responsible for maintaining the confidentiality of your account and password and for restricting access to your computer, and you agree to accept responsibility for all activities that occur under your account or password. You may not assign or otherwise transfer your account to any other person or entity. You acknowledge that Aether Technologies Inc. is not responsible for third party access to your account that results from theft or misappropriation of your account. Aether Technologies Inc. and its associates reserve the right to refuse or cancel service, terminate accounts, or remove or edit content in our sole discretion.

Children Under Thirteen

Aether Technologies Inc. does not knowingly collect, either online or offline, personal information from persons under the age of thirteen. If you are under 18, you may use Getaether.net or Aether App only with permission of a parent or guardian.

Links to Third Party Sites/Third Party Services

Getaether.net or Aether App may contain links to other websites ("Linked Sites"). The Linked Sites are not under the control of Aether Technologies Inc. and Aether Technologies Inc. is not responsible for the contents of any Linked Site, including without limitation any link contained in a Linked Site, or any changes or updates to a Linked Site. Aether Technologies Inc. is providing these links to you only as a convenience, and the inclusion of any link does not imply endorsement by Aether Technologies Inc. of the site or any association with its operators.

Certain services made available via Getaether.net or Aether App are delivered by third party sites and organizations. By using any product, service or functionality originating from the Getaether.net domain or from the Aether App, you hereby acknowledge and consent that Aether Technologies Inc. may share such information and data with any third party with whom Aether Technologies Inc. has a contractual relationship to provide the requested product, service or functionality on behalf of Getaether.net or Aether App users and customers.

No Unlawful or Prohibited Use/Intellectual Property

You are granted a non-exclusive, non-transferable, revocable license to access and use Getaether.net and Aether App strictly in accordance with these terms of use. As a condition of your use of the Site and App, you warrant to Aether Technologies Inc. that you will not use the Site or App for any purpose that is unlawful or prohibited by these Terms. You may not use the Site or App in any manner which could damage, disable, overburden, or impair the Site or App or interfere with any other party's use and enjoyment of the Site or App. You may not obtain or attempt to obtain any materials or information through any means not intentionally made available or provided for through the Site or App.

All content included as part of the Service, such as text, graphics, logos, images, as well as the compilation thereof, and any software used on the Site or App, is the property of Aether Technologies Inc. or its suppliers and protected by copyright and other laws that protect intellectual property and proprietary rights. You agree to observe and abide by all copyright and other proprietary notices, legends or other restrictions contained in any such content and will not make any changes thereto.

You will not modify, publish, transmit, reverse engineer, participate in the transfer or sale, create derivative works, or in any way exploit any of the content, in whole or in part, found on the Site or App. Aether Technologies Inc. content is not for resale. Your use of the Site or App does not entitle you to make any unauthorized use of any protected content, and in particular you will not delete or alter any proprietary rights or attribution notices in any content. You will use protected content solely for your personal use, and will make no other use of the content without the express written permission of Aether Technologies Inc. and the copyright owner. You agree that you do not acquire any ownership rights in any protected content. We do not grant you any licenses, express or implied, to the intellectual property of Aether Technologies Inc. or our licensors except as expressly authorized by these Terms.

Use of Communication Services

The Site or App may contain bulletin board services, chat areas, news groups, forums, communities, personal web pages, calendars, and/or other message or communication facilities designed to enable you to communicate with the public at large or with a group (collectively, "Communication Services"). You agree to use the Communication Services only to post, send and receive messages and material that are proper and related to the particular Communication Service.

By way of example, and not as a limitation, you agree that when using a Communication Service, you will not: defame, abuse, harass, stalk, threaten or otherwise violate the legal rights (such as rights of privacy and publicity) of others; publish, post, upload, distribute or disseminate any inappropriate, profane, defamatory, infringing, obscene, indecent or unlawful topic, name, material or information; upload files that contain software or other material protected by intellectual property laws (or by rights of privacy of publicity) unless you own or control the rights thereto or have received all necessary consents; upload files that contain viruses, corrupted files, or any other similar software or programs that may damage the operation of another's computer; advertise or offer to sell or buy any goods or services for any business purpose, unless such Communication Service specifically allows such messages; conduct or forward surveys, contests, pyramid schemes or chain letters; download any file posted by another user of a Communication Service that you know, or reasonably should know, cannot be legally distributed in such manner; falsify or delete any author attributions, legal or other proper notices or proprietary designations or labels of the origin or source of software or other material contained in a file that is uploaded; restrict or inhibit any other user from using and enjoying the Communication Services; violate any code of conduct or other guidelines which may be applicable for any particular Communication Service; harvest or otherwise collect information about others, including e-mail addresses, without their consent; violate any applicable laws or regulations.

Aether Technologies Inc. has no obligation to monitor the Communication Services. However, Aether Technologies Inc. reserves the right to review materials posted to a Communication Service and to remove any materials in its sole discretion. Aether Technologies Inc. reserves the right to terminate your access to any or all of the Communication Services at any time without notice for any reason whatsoever.

Aether Technologies Inc. reserves the right at all times to disclose any information as necessary to satisfy any applicable law, regulation, legal process or governmental request, or to edit, refuse to post or to remove any information or materials, in whole or in part, in Aether Technologies Inc.'s sole discretion.

Always use caution when giving out any personally identifying information about yourself or your children in any Communication Service. Aether Technologies Inc. does not control or endorse the content, messages or information found in any Communication Service and, therefore, Aether Technologies Inc. specifically disclaims any liability with regard to the Communication Services and any actions resulting from your participation in any Communication Service. Managers and hosts are not authorized Aether Technologies Inc. spokespersons, and their views do not necessarily reflect those of Aether Technologies Inc..

Materials uploaded to a Communication Service may be subject to posted limitations on usage, reproduction and/or dissemination. You are responsible for adhering to such limitations if you upload the materials.

Peer-to-Peer Nature

Aether App is a peer-to-peer app. Using Aether means that your computer or device shares the content created by users of Aether and other compatible apps available on the Aether network to other participants in the network. It is your responsibility to monitor whether this content shared from your computer or device is in compliance with your local jurisdiction.

Content Policy

Aether tries to provide a lot of flexibility in what kind of content it can be used for, however, there are certain kinds of content that are legally problematic. Following types of content are not allowed:

* Illegal content
* Involuntary pornography
* Sexual or suggestive content involving minors
* Content that encourages or incites violence
* Content that threatens, or encourages others to do so
* Content that is personal and confidential information
* Content that uses Aether to solicit or facilitate any transaction or gift involving certain goods and services
* Content that is spam

Illegal content

Do not use Aether for illegal content.

Involuntary pornography

Do not disseminate images, video, depicting any person in any act of sexual conduct, apparently created or published without their explicit permission.

This applies even if a person is fully clothed, and/or in public. Images of body parts recontextualised in a salacious manner (e.g. 'Upskirt', 'Creepshots') also fall under this restriction.

Sexual or suggestive content involving minors

Aether does not allow any sexual or suggestive content, involving minors, or person(s) who appear to be minors, whether be it videos, photos, 3D renderings, animation, or in any other form.

Content that encourages or incites violence

Do not post content that encourages violence or physical harm against other people, against self, or against animals.

Content that threatens or encourages others to do so

Do not post content that can make a reasonable person conclude that he or she is under the threat of physical harm.

Content that is personal and confidential information

Don't post someone else's personal information online. (e.g. 'Doxxing') Public figures can sometimes be an exception to this rule depending on the context. As an example, if you're posting the contact email, of your local representative, that is fine, if you're posting his or her home address, or garage PIN code, it is not.

Content that uses Aether to solicit or facilitate any transaction or gift involving certain goods and services


Aether is not intended to be a marketplace.


Specifically, using Aether to seek, buy, sell, or otherwise facilitate transactions of these types of items is not allowed:

- Firearms, ammunition or explosives of any kind, or instructions to make them, whether be it written, in computer code (e.g. 3D printing files)

- Drugs, Alcohol, Tobacco or Tobacco products, whether be it of organic or synthetic origin. (e.g. e-cigarettes)

- Paid services involving physical sexual contact (e.g. prostitution)

- Stolen goods

- Personal information

- Counterfeit items, documents or currency

- Scams, and otherwise fraudulent services

Content that is spam

Don't spam Aether.

Behavior Policy

Do not engage in vote manipulation

Do not create multiple accounts with the purpose of being able to vote multiple times on the same content or in the same moderator election cycle.

Creating multiple accounts to post in different communities is okay, so long as you don't vote on content or people twice.

Do not break Aether

Aether is a commons that works with the collaboration and participation of all its users. Don't break it, or do anything that interferes with its normal use. By way of example, and not as a limitation, do not post base64 encoded images (or similar) in text form.

White-hat penetration testing and security research are exceptions to this rule, as long as responsible disclosure policies are followed.

Materials Provided to Getaether.net, Aether App, or Posted on Any Aether Technologies Inc. Web Page

Aether Technologies Inc. does not claim ownership of the materials you provide to Getaether.net (including feedback and suggestions), the Aether App, or post, upload, input or submit to any Aether Technologies Inc. Site or our associated services (collectively "Submissions"). However, by posting, uploading, inputting, providing or submitting your Submission you are granting Aether Technologies Inc., our affiliated companies and necessary sublicensees permission to use your Submission in connection with the operation of their Internet businesses including, without limitation, the rights to: copy, distribute, transmit, publicly display, publicly perform, reproduce, edit, translate and reformat your Submission; and to publish your name in connection with your Submission.

In addition, by posting, uploading, inputting, providing or submitting your Submission, you are perpetually, and nonrevocably licensing your Submission under Creative Commons BY-SA license. This license can be found at the following URL. https://creativecommons.org/licenses/by-sa/4.0/

No compensation will be paid with respect to the use of your Submission, as provided herein. Aether Technologies Inc. is under no obligation to post or use any Submission you may provide and may remove any Submission at any time in Aether Technologies Inc.'s sole discretion.

By posting, uploading, inputting, providing or submitting your Submission you warrant and represent that you own or otherwise control all of the rights to your Submission as described in this section including, without limitation, all the rights necessary for you to provide, post, upload, input or submit the Submissions.

Third Party Accounts

You will be able to connect your Aether Technologies Inc. account to third party accounts. By connecting your Aether Technologies Inc. account to your third party account, you acknowledge and agree that you are consenting to the continuous release of information about you to others (in accordance with your privacy settings on those third party sites). If you do not want information about you to be shared in this manner, do not use this feature.

International Users

The Service is controlled, operated and administered by Aether Technologies Inc. from our offices within the USA. If you access the Service from a location outside the USA, you are responsible for compliance with all local laws. You agree that you will not use the Aether Technologies Inc. Content accessed through Getaether.net and Aether App in any country or in any manner prohibited by any applicable laws, restrictions or regulations.

Indemnification

You agree to indemnify, defend and hold harmless Aether Technologies Inc., its officers, directors, employees, agents and third parties, for any losses, costs, liabilities and expenses (including reasonable attorney's fees) relating to or arising out of your use of or inability to use the Site, the App, or services, any user postings made by you, your violation of any terms of this Agreement or your violation of any rights of a third party, or your violation of any applicable laws, rules or regulations. Aether Technologies Inc. reserves the right, at its own cost, to assume the exclusive defense and control of any matter otherwise subject to indemnification by you, in which event you will fully cooperate with Aether Technologies Inc. in asserting any available defenses.

Arbitration

In the event the parties are not able to resolve any dispute between them arising out of or concerning these Terms and Conditions, or any provisions hereof, whether in contract, tort, or otherwise at law or in equity for damages or any other relief, then such dispute shall be resolved only by final and binding arbitration pursuant to the Federal Arbitration Act, conducted by a single neutral arbitrator and administered by the American Arbitration Association, or a similar arbitration service selected by the parties, in a location mutually agreed upon by the parties. The arbitrator's award shall be final, and judgment may be entered upon it in any court having jurisdiction. In the event that any legal or equitable action, proceeding or arbitration arises out of or concerns these Terms and Conditions, the prevailing party shall be entitled to recover its costs and reasonable attorney's fees. The parties agree to arbitrate all disputes and claims in regards to these Terms and Conditions or any disputes arising as a result of these Terms and Conditions, whether directly or indirectly, including Tort claims that are a result of these Terms and Conditions. The parties agree that the Federal Arbitration Act governs the interpretation and enforcement of this provision. The entire dispute, including the scope and enforceability of this arbitration provision shall be determined by the Arbitrator. This arbitration provision shall survive the termination of these Terms and Conditions.

Class Action Waiver

Any arbitration under these Terms and Conditions will take place on an individual basis; class arbitrations and class/representative/collective actions are not permitted. THE PARTIES AGREE THAT A PARTY MAY BRING CLAIMS AGAINST THE OTHER ONLY IN EACH'S INDIVIDUAL CAPACITY, AND NOT AS A PLAINTIFF OR CLASS MEMBER IN ANY PUTATIVE CLASS, COLLECTIVE AND/ OR REPRESENTATIVE PROCEEDING, SUCH AS IN THE FORM OF A PRIVATE ATTORNEY GENERAL ACTION AGAINST THE OTHER. Further, unless both you and Aether Technologies Inc. agree otherwise, the arbitrator may not consolidate more than one person's claims, and may not otherwise preside over any form of a representative or class proceeding.

Liability Disclaimer

THE INFORMATION, SOFTWARE, PRODUCTS, AND SERVICES INCLUDED IN OR AVAILABLE THROUGH THE SITE OR APP MAY INCLUDE INACCURACIES OR TYPOGRAPHICAL ERRORS. CHANGES ARE PERIODICALLY ADDED TO THE INFORMATION HEREIN. AETHER TECHNOLOGIES INC. AND/OR ITS SUPPLIERS MAY MAKE IMPROVEMENTS AND/OR CHANGES IN THE SITE OR APP AT ANY TIME.

AETHER TECHNOLOGIES INC. AND/OR ITS SUPPLIERS MAKE NO REPRESENTATIONS ABOUT THE SUITABILITY, RELIABILITY, AVAILABILITY, TIMELINESS, AND ACCURACY OF THE INFORMATION, SOFTWARE, PRODUCTS, SERVICES AND RELATED GRAPHICS CONTAINED ON THE SITE OR APP FOR ANY PURPOSE. TO THE MAXIMUM EXTENT PERMITTED BY APPLICABLE LAW, ALL SUCH INFORMATION, SOFTWARE, PRODUCTS, SERVICES AND RELATED GRAPHICS ARE PROVIDED "AS IS" WITHOUT WARRANTY OR CONDITION OF ANY KIND. AETHER TECHNOLOGIES INC. AND/OR ITS SUPPLIERS HEREBY DISCLAIM ALL WARRANTIES AND CONDITIONS WITH REGARD TO THIS INFORMATION, SOFTWARE, PRODUCTS, SERVICES AND RELATED GRAPHICS, INCLUDING ALL IMPLIED WARRANTIES OR CONDITIONS OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, TITLE AND NON-INFRINGEMENT.

TO THE MAXIMUM EXTENT PERMITTED BY APPLICABLE LAW, IN NO EVENT SHALL AETHER TECHNOLOGIES INC. AND/OR ITS SUPPLIERS BE LIABLE FOR ANY DIRECT, INDIRECT, PUNITIVE, INCIDENTAL, SPECIAL, CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER INCLUDING, WITHOUT LIMITATION, DAMAGES FOR LOSS OF USE, DATA OR PROFITS, ARISING OUT OF OR IN ANY WAY CONNECTED WITH THE USE OR PERFORMANCE OF THE SITE OR APP, WITH THE DELAY OR INABILITY TO USE THE SITE, APP OR RELATED SERVICES, THE PROVISION OF OR FAILURE TO PROVIDE SERVICES, OR FOR ANY INFORMATION, SOFTWARE, PRODUCTS, SERVICES AND RELATED GRAPHICS OBTAINED THROUGH THE SITE, APP, OR OTHERWISE ARISING OUT OF THE USE OF THE SITE OR APP, WHETHER BASED ON CONTRACT, TORT, NEGLIGENCE, STRICT LIABILITY OR OTHERWISE, EVEN IF AETHER TECHNOLOGIES INC. OR ANY OF ITS SUPPLIERS HAS BEEN ADVISED OF THE POSSIBILITY OF DAMAGES. BECAUSE SOME STATES/JURISDICTIONS DO NOT ALLOW THE EXCLUSION OR LIMITATION OF LIABILITY FOR CONSEQUENTIAL OR INCIDENTAL DAMAGES, THE ABOVE LIMITATION MAY NOT APPLY TO YOU. IF YOU ARE DISSATISFIED WITH ANY PORTION OF THE SITE OR APP, OR WITH ANY OF THESE TERMS OF USE, YOUR SOLE AND EXCLUSIVE REMEDY IS TO DISCONTINUE USING THE SITE OR APP.

Termination/Access Restriction

Aether Technologies Inc. reserves the right, in its sole discretion, to terminate your access to the Site or App and the related services or any portion thereof at any time, without notice. To the maximum extent permitted by law, this agreement is governed by the laws of the State of California and you hereby consent to the exclusive jurisdiction and venue of courts in California in all disputes arising out of or relating to the use of the Site or App. Use of the Site or App is unauthorized in any jurisdiction that does not give effect to all provisions of these Terms, including, without limitation, this section.

You agree that no joint venture, partnership, employment, or agency relationship exists between you and Aether Technologies Inc. as a result of this agreement or use of the Site or App. Aether Technologies Inc.'s performance of this agreement is subject to existing laws and legal process, and nothing contained in this agreement is in derogation of Aether Technologies Inc.'s right to comply with governmental, court and law enforcement requests or requirements relating to your use of the Site and App or information provided to or gathered by Aether Technologies Inc. with respect to such use. If any part of this agreement is determined to be invalid or unenforceable pursuant to applicable law including, but not limited to, the warranty disclaimers and liability limitations set forth above, then the invalid or unenforceable provision will be deemed superseded by a valid, enforceable provision that most closely matches the intent of the original provision and the remainder of the agreement shall continue in effect.

Unless otherwise specified herein, this agreement constitutes the entire agreement between the user and Aether Technologies Inc. with respect to the Site and App and it supersedes all prior or contemporaneous communications and proposals, whether electronic, oral or written, between the user and Aether Technologies Inc. with respect to the Site and App. A printed version of this agreement and of any notice given in electronic form shall be admissible in judicial or administrative proceedings based upon or relating to this agreement to the same extent and subject to the same conditions as other business documents and records originally generated and maintained in printed form. It is the express wish to the parties that this agreement and all related documents be written in English.

Changes to Terms

Aether Technologies Inc. reserves the right, in its sole discretion, to change the Terms under which Getaether.net and Aether App is offered. The most current version of the Terms will supersede all previous versions. Aether Technologies Inc. encourages you to periodically review the Terms to stay informed of our updates.

Contact Us

Aether Technologies, Inc.

1 Embarcadero Ctr Ste SL12 #2025

San Francisco, CA 94111, United States


Email Address:

legal@getaether.net
`
</script>

<style lang="scss" scoped>
@import '../../../scss/globals';
// .onboard-sublocation {
//   display: flex;
//   flex: 1;
//   min-height: 100%;
//   .onboard-carrier {
//     margin: auto;
//     display: flex;
//     width: 518px;
//     flex-direction: column;
//     position: relative;
//     bottom: 5%;
//   }
// }

.tocbox {
  height: 157px;
  overflow-y: scroll;
  // @include generateScrollbar($dark-base);
  // @include generateScrollbar($a-grey-100)
  &::-webkit-scrollbar-track {
    background-color: $a-transparent;
  }
  &::-webkit-scrollbar-thumb {
    background-color: #fff;
  }
  background-color: rgba(0, 0, 0, 0.3);
  font-size: 80%;
  padding: 15px;
  border-radius: 5px;
}

.continue-box {
  display: flex;
  margin-top: 15px;
  .next-button {
    @extend %link-hover-ghost-extenders-disable;
    background-color: $a-transparent;
    color: $a-grey-800;
    &:hover {
      background-color: $a-grey-800;
      color: $mid-base;
    }
    animation-duration: 1s;
    animation-name: DELAY_VISIBLE;
    margin-right: 10px;
  }
}

@keyframes DELAY_VISIBLE {
  0% {
    opacity: 0;
  }
  60% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition-property: opacity;
  transition-duration: 0.25s;
}

.fade-enter-active {
  transition-delay: 0.25s;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}
</style>

<style lang="scss">
.onboard-markdown {
  font-size: 120%;
  p {
    font-family: 'SSP Regular';
  }
}
</style>
