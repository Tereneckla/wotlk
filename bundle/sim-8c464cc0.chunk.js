import{A as e,eq as t,er as a,dY as n,L as l,m as s,n as i,es as o,et as d,eu as r,ev as c,ew as p,ex as m,ey as u,E as h,co as I,cj as g,cm as S,bD as f,aR as v,dU as b,K as y,T as A,a2 as T,F as w,bd as O,be as M,aD as P,bn as R,w as C,B as W,aE as L}from"./detailed_results-7b150079.chunk.js";import{m as N,u as k,a as F,b as B,c as E,B as D,I as G,T as x}from"./individual_sim_ui-7ca50b32.chunk.js";import{T as H}from"./totem_inputs-27ceb772.chunk.js";const J=N({fieldName:"bloodlust",id:e.fromSpellId(2825)}),j=k({fieldName:"shield",values:[{value:t.NoShield,tooltip:"No Shield"},{actionId:e.fromSpellId(33736),value:t.WaterShield},{actionId:e.fromSpellId(25472),value:t.LightningShield}]}),q={inputs:[F({fieldName:"type",label:"Type",values:[{name:"Adaptive",value:a.Adaptive,tooltip:"Dynamically adapts based on available mana to maximize CL casts without going OOM."},{name:"Manual",value:a.Manual,tooltip:"Allows custom selection of which spells to use and to modify cast conditions."}]}),B({fieldName:"inThunderstormRange",label:"In Thunderstorm Range",labelTooltip:"Thunderstorm will hit all targets when cast. Ignores knockback.",showWhen:e=>e.getTalents().thunderstorm}),E({fieldName:"lvbFsWaitMs",label:"Max wait for LvB/FS (ms)",labelTooltip:"Amount of time the sim will wait if FS is about to fall off or LvB CD is about to come up. Setting to 0 will default to 175ms"}),B({fieldName:"useChainLightning",label:"Use Chain Lightning in Rotation",labelTooltip:"Use Chain Lightning in rotation",enableWhen:e=>e.getRotation().type==a.Manual}),B({fieldName:"useClOnlyGap",label:"Use CL only as gap filler",labelTooltip:"Use CL to fill short gaps in LvB CD instead of on CD.",enableWhen:e=>e.getRotation().type==a.Manual&&e.getRotation().useChainLightning}),E({fieldName:"clMinManaPer",label:"Min mana percent to use Chain Lightning",labelTooltip:"Customize minimum mana level to cast Chain Lightning. 0 will spam until OOM.",enableWhen:e=>e.getRotation().type==a.Manual&&e.getRotation().useChainLightning}),B({fieldName:"useFireNova",label:"Use Fire Nova in Rotation",labelTooltip:"Fire Nova will hit all targets when cast.",enableWhen:e=>e.getRotation().type==a.Manual}),E({fieldName:"fnMinManaPer",label:"Min mana percent to use FireNova",labelTooltip:"Customize minimum mana level to cast Fire Nova. 0 will spam until OOM.",enableWhen:e=>e.getRotation().type==a.Manual&&e.getRotation().useFireNova}),B({fieldName:"overwriteFlameshock",label:"Allow Flameshock to be overwritten",labelTooltip:"Will use flameshock at the end of the duration even if its still ticking if there isn't enough time to cast lavaburst before expiring.",enableWhen:e=>e.getRotation().type==a.Manual}),B({fieldName:"alwaysCritLvb",label:"Only cast Lavaburst with FS",labelTooltip:"Will only cast Lavaburst if Flameshock will be active when the cast finishes.",enableWhen:e=>e.getRotation().type==a.Manual}),B({fieldName:"useThunderstorm",label:"Allow Thunderstorm to be cast.",labelTooltip:"Disabling this will stop thunderstorm from being cast entirely.",enableWhen:e=>e.getRotation().type==a.Manual,showWhen:e=>e.getTalents().thunderstorm})]},U={name:"Standard",data:n.create({talentsString:"0532001523212351322301051-00504",glyphs:l.create({major1:s.GlyphOfLightningBolt,major2:s.GlyphOfTotemOfWrath,major3:s.ShamanMajorGlyphNone,minor1:i.GlyphOfThunderstorm,minor2:i.GlyphOfWaterShield,minor3:i.GlyphOfGhostWolf})})},V=o.create({totems:d.create({earth:r.StrengthOfEarthTotem,air:c.WrathOfAirTotem,fire:p.TotemOfWrath,water:m.ManaSpringTotem,useFireElemental:!0}),type:a.Adaptive,fnMinManaPer:66,clMinManaPer:33,useChainLightning:!1,useFireNova:!1,useThunderstorm:!0}),_=u.create({shield:t.WaterShield,bloodlust:!0}),z=h.create({defaultPotion:I.HastePotion,flask:g.FlaskOfBlindingLight,food:S.FoodBlackenedBasilisk}),K={name:"Pre-raid Preset",tooltip:D,gear:f.fromJsonString('{"items": [\n\t\t{"id":37180,"enchant":3820,"gems":[41285,42144]},\n\t\t{"id":37595},\n\t\t{"id":37673,"enchant":3810,"gems":[42144]},\n\t\t{"id":41610,"enchant":3722},\n\t\t{"id":39592,"enchant":3832,"gems":[42144,40025]},\n\t\t{"id":37788,"enchant":2332,"gems":[0]},\n\t\t{"id":39593,"enchant":3246,"gems":[40051,0]},\n\t\t{"id":40696,"gems":[40049,39998]},\n\t\t{"id":37791,"enchant":3719},\n\t\t{"id":44202,"enchant":3826,"gems":[39998]},\n\t\t{"id":43253,"gems":[40027]},\n\t\t{"id":37694},\n\t\t{"id":40682},\n\t\t{"id":37873},\n\t\t{"id":41384,"enchant":3834},\n\t\t{"id":40698},\n\t\t{"id":40708}\n  ]}')},Y={name:"P1 Preset",tooltip:D,gear:f.fromJsonString('{"items": [\n\t\t{"id":40516,"enchant":3820,"gems":[41285,40027]},\n\t\t{"id":44661,"gems":[39998]},\n\t\t{"id":40286,"enchant":3810},\n\t\t{"id":44005,"enchant":3722,"gems":[40027]},\n\t\t{"id":40514,"enchant":3832,"gems":[42144,42144]},\n\t\t{"id":40324,"enchant":2332,"gems":[42144,0]},\n\t\t{"id":40302,"enchant":3246,"gems":[0]},\n\t\t{"id":40301,"gems":[40014]},\n\t\t{"id":40560,"enchant":3721},\n\t\t{"id":40519,"enchant":3826},\n\t\t{"id":37694},\n\t\t{"id":40399},\n\t\t{"id":40432},\n\t\t{"id":40255},\n\t\t{"id":40395,"enchant":3834},\n\t\t{"id":40401,"enchant":1128},\n\t\t{"id":40267}\n  ]}')},Q={name:"P2 Preset",tooltip:D,gear:f.fromJsonString('{"items": [\n        {"id":46209,"enchant":3820,"gems":[41285,40048]},\n        {"id":45933,"gems":[39998]},\n        {"id":46211,"enchant":3810,"gems":[39998]},\n        {"id":45242,"enchant":3722,"gems":[39998]},\n        {"id":46206,"enchant":3832,"gems":[39998,39998]},\n        {"id":45460,"enchant":2332,"gems":[39998,0]},\n        {"id":45665,"enchant":3604,"gems":[39998,39998,0]},\n        {"id":45616,"enchant":3599,"gems":[39998,39998,39998]},\n        {"id":46210,"enchant":3721,"gems":[39998,40027]},\n        {"id":45537,"enchant":3606,"gems":[39998,40027]},\n        {"id":46046,"gems":[39998]},\n        {"id":45495,"gems":[39998]},\n        {"id":45518},\n        {"id":40255},\n        {"id":45612,"enchant":3834,"gems":[39998]},\n        {"id":45470,"enchant":1128,"gems":[39998]},\n        {"id":40267}\n      ]}')},X={name:"P3 Preset [H]",enableWhen:e=>e.getFaction()==v.Horde,tooltip:D,gear:f.fromJsonString('{"items": [\n        {"id":48328,"enchant":3820,"gems":[41285,40153]},\n        {"id":47468,"gems":[40155]},\n        {"id":48330,"enchant":3810,"gems":[40113]},\n        {"id":47551,"enchant":3722,"gems":[40113]},\n        {"id":48326,"enchant":3832,"gems":[40113,40132]},\n        {"id":45460,"enchant":2332,"gems":[40113,0]},\n        {"id":48327,"enchant":3604,"gems":[40155,0]},\n        {"id":47447,"enchant":3599,"gems":[40132,40113,40113]},\n        {"id":47479,"enchant":3721,"gems":[40113,40113,40113]},\n        {"id":47456,"enchant":3606,"gems":[40113,40113]},\n        {"id":46046,"gems":[40155]},\n        {"id":45495,"gems":[40113]},\n        {"id":47477},\n        {"id":45518},\n        {"id":47422,"enchant":3834,"gems":[40113]},\n        {"id":47448,"enchant":1128,"gems":[40155]},\n        {"id":47666}\n      ]\n    }')},Z={name:"P3 Preset [A]",enableWhen:e=>e.getFaction()==v.Alliance,tooltip:D,gear:f.fromJsonString('{"items": [\n        {"id":48323,"enchant":3820,"gems":[41285,40155]},\n        {"id":47144,"gems":[40155]},\n        {"id":48321,"enchant":3810,"gems":[40113]},\n        {"id":47552,"enchant":3722,"gems":[40113]},\n        {"id":48325,"enchant":3832,"gems":[40113,40132]},\n        {"id":45460,"enchant":2332,"gems":[40113,0]},\n        {"id":48324,"enchant":3604,"gems":[40155,0]},\n        {"id":47084,"enchant":3599,"gems":[40132,40113,40113]},\n        {"id":47190,"enchant":3721,"gems":[40113,40113,40113]},\n        {"id":47099,"enchant":3606,"gems":[40113,40113]},\n        {"id":46046,"gems":[40155]},\n        {"id":45495,"gems":[40113]},\n        {"id":47188},\n        {"id":45518},\n        {"id":46980,"enchant":3834,"gems":[40113]},\n        {"id":47085,"enchant":1128,"gems":[40155]},\n        {"id":47666}\n      ]\n    }')},$={name:"Legacy",rotation:b.create({specRotationOptionsJson:o.toJsonString(V)})},ee={name:"Basic APL",rotation:b.create({specRotationOptionsJson:o.toJsonString(V),rotation:y.fromJsonString('{\n      "type": "TypeAPL",\n      "prepullActions": [\n\t\t\t  {"action":{"castSpell":{"spellId":{"spellId":3738}}},"doAtValue":{"const":{"val":"-6s"}}},\n\t\t\t  {"action":{"castSpell":{"spellId":{"spellId":58643}}},"doAtValue":{"const":{"val":"-5s"}}},\n\t\t\t  {"action":{"castSpell":{"spellId":{"spellId":58774}}},"doAtValue":{"const":{"val":"-4s"}}},\n\t\t\t  {"action":{"castSpell":{"spellId":{"spellId":57722}}},"doAtValue":{"const":{"val":"-3s"}}},\n\t\t\t  {"action":{"castSpell":{"spellId":{"spellId":58704}}},"doAtValue":{"const":{"val":"-2s"}}},\n\t\t\t  {"action":{"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}},"doAtValue":{"const":{"val":"-1s"}}}\n      ],\n      "priorityList": [\n        {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"2s"}}}},{"spellIsReady":{"spellId":{"tag":-1,"spellId":2825}}}]}},"castSpell":{"spellId":{"tag":-1,"spellId":2825}}}},\n        {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"2s"}}}},{"spellIsReady":{"spellId":{"spellId":2825}}}]}},"castSpell":{"spellId":{"spellId":2825}}}},\n        {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":26297}}},{"spellIsReady":{"spellId":{"spellId":16166}}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":26297}}},{"castSpell":{"spellId":{"spellId":16166}}}]}}},\n        {"action":{"condition":{"and":{"vals":[{"not":{"val":{"spellIsReady":{"spellId":{"spellId":26297}}}}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":16166}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":64701}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":26297}}}}}]}},"castSpell":{"spellId":{"spellId":54758}}}},\n        {"action":{"condition":{"and":{"vals":[{"spellIsReady":{"spellId":{"spellId":2894}}},{"or":{"vals":[{"auraIsActive":{"auraId":{"itemId":40255}}},{"auraIsActive":{"auraId":{"itemId":40682}}},{"auraIsActive":{"auraId":{"itemId":37660}}},{"auraIsActive":{"auraId":{"itemId":45518}}},{"auraIsActive":{"auraId":{"itemId":54572}}},{"auraIsActive":{"auraId":{"itemId":54588}}},{"auraIsActive":{"auraId":{"itemId":47213}}},{"auraIsActive":{"auraId":{"itemId":45490}}},{"auraIsActive":{"auraId":{"itemId":50348}}},{"auraIsActive":{"auraId":{"itemId":50353}}},{"auraIsActive":{"auraId":{"itemId":50360}}},{"auraIsActive":{"auraId":{"itemId":50365}}},{"auraIsActive":{"auraId":{"itemId":50345}}},{"auraIsActive":{"auraId":{"itemId":50340}}},{"auraIsActive":{"auraId":{"itemId":50398}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"itemId":45308}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"itemId":40432}}},"rhs":{"const":{"val":"10"}}}},{"auraIsActive":{"auraId":{"spellId":55637}}}]}}]}},"strictSequence":{"actions":[{"castSpell":{"spellId":{"spellId":33697}}},{"castSpell":{"spellId":{"itemId":40212}}},{"castSpell":{"spellId":{"itemId":37873}}},{"castSpell":{"spellId":{"itemId":45148}}},{"castSpell":{"spellId":{"itemId":48724}}},{"castSpell":{"spellId":{"itemId":50357}}},{"castSpell":{"spellId":{"spellId":2894}}}]}}},\n        {"action":{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":2894}}}}},{"not":{"val":{"dotIsActive":{"spellId":{"spellId":58704}}}}}]}},"castSpell":{"spellId":{"spellId":58704}}}},\n        {"action":{"multidot":{"spellId":{"spellId":49233},"maxDots":3,"maxOverlap":{"const":{"val":"0ms"}}}}},\n        {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"numberTargets":{}},"rhs":{"const":{"val":"1"}}}},{"spellIsReady":{"spellId":{"spellId":49271}}}]}},"castSpell":{"spellId":{"spellId":49271}}}},\n        {"action":{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"dotRemainingTime":{"spellId":{"spellId":49233}}},"rhs":{"const":{"val":"2"}}}}]}},"castSpell":{"spellId":{"spellId":60043}}}},\n        {"action":{"castSpell":{"spellId":{"spellId":49238}}}}\n      ]\n    }')})};class te extends G{constructor(e,t){super(e,t,{cssClass:"elemental-shaman-sim-ui",cssScheme:"shaman",knownIssues:[],warnings:[e=>({updateOn:A.onAny([e.player.rotationChangeEmitter,e.player.currentStatsEmitter]),getContent:()=>{const t=e.player.getCurrentStats().sets.includes("Skyshatter Regalia (2pc)"),a=e.player.getSpecOptions().totems,n=a&&a.earth&&a.air&&a.fire&&a.water;return t&&!n?"T6 2pc bonus is equipped, but inactive because not all 4 totem types are being used.":""}})],epStats:[T.StatIntellect,T.StatSpellPower,T.StatSpellHit,T.StatSpellCrit,T.StatSpellHaste,T.StatMP5],epReferenceStat:T.StatSpellPower,displayStats:[T.StatHealth,T.StatMana,T.StatStamina,T.StatIntellect,T.StatSpellPower,T.StatSpellHit,T.StatSpellCrit,T.StatSpellHaste,T.StatMP5],modifyDisplayStats:e=>{let t=new w;return t=t.addStat(T.StatSpellHit,e.getTalents().elementalPrecision*O),t=t.addStat(T.StatSpellCrit,1*e.getTalents().tidalMastery*M),{talents:t}},defaults:{gear:Y.gear,epWeights:w.fromMap({[T.StatIntellect]:.22,[T.StatSpellPower]:1,[T.StatSpellCrit]:.67,[T.StatSpellHaste]:1.29,[T.StatMP5]:.08}),consumes:z,rotation:V,talents:U.data,specOptions:_,raidBuffs:P.create({arcaneBrilliance:!0,divineSpirit:!0,giftOfTheWild:R.TristateEffectImproved,moonkinAura:R.TristateEffectImproved,sanctifiedRetribution:!0}),partyBuffs:C.create({}),individualBuffs:W.create({blessingOfKings:!0,blessingOfWisdom:2,vampiricTouch:!0}),debuffs:L.create({faerieFire:R.TristateEffectImproved,judgementOfWisdom:!0,misery:!0,curseOfElements:!0,shadowMastery:!0})},playerIconInputs:[j,J],rotationInputs:q,includeBuffDebuffInputs:[],excludeBuffDebuffInputs:[],otherInputs:{inputs:[x]},customSections:[H],encounterPicker:{showExecuteProportion:!1},presets:{talents:[U],rotations:[ee,$],gear:[K,Y,Q,Z,X]}})}}export{V as D,te as E,Y as P,U as S,_ as a,z as b,Q as c,Z as d,X as e};
//# sourceMappingURL=sim-8c464cc0.chunk.js.map
